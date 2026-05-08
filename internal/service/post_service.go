package service

import (
	"errors"
	"time"

	"github.com/qiblatdigital/zoztool-api/internal/model"
	"github.com/qiblatdigital/zoztool-api/internal/pkg/threads"
	"github.com/qiblatdigital/zoztool-api/internal/repository"
)

type PostService struct {
	postRepo          *repository.PostRepository
	postTargetRepo    *repository.PostTargetRepository
	socialAccountRepo *repository.SocialAccountRepository
	threadsClient     *threads.Client
}

func NewPostService(
	postRepo *repository.PostRepository,
	postTargetRepo *repository.PostTargetRepository,
	socialAccountRepo *repository.SocialAccountRepository,
	threadsClient *threads.Client,
) *PostService {
	return &PostService{
		postRepo:          postRepo,
		postTargetRepo:    postTargetRepo,
		socialAccountRepo: socialAccountRepo,
		threadsClient:     threadsClient,
	}
}

type CreatePostRequest struct {
	Title            string     `json:"title"`
	Content          string     `json:"content" binding:"required"`
	ScheduledAt      *time.Time `json:"scheduled_at"`
	TargetAccountIDs []string   `json:"target_account_ids" binding:"required,min=1"`
}

type UpdatePostRequest struct {
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	ScheduledAt *time.Time `json:"scheduled_at"`
}

type PostWithTargets struct {
	model.Post
	Targets []model.PostTarget `json:"targets"`
}

func (s *PostService) Create(userID string, req CreatePostRequest) (*PostWithTargets, error) {
	status := model.PostStatusDraft
	if req.ScheduledAt != nil {
		status = model.PostStatusScheduled
	}

	post := &model.Post{
		UserID:      userID,
		Title:       req.Title,
		Content:     req.Content,
		Status:      status,
		ScheduledAt: req.ScheduledAt,
	}
	if err := s.postRepo.Create(post); err != nil {
		return nil, errors.New("failed to create post")
	}

	targets := make([]model.PostTarget, 0, len(req.TargetAccountIDs))
	for _, accountID := range req.TargetAccountIDs {
		targets = append(targets, model.PostTarget{
			PostID:          post.ID,
			SocialAccountID: accountID,
			Status:          model.PostTargetStatusPending,
		})
	}
	if err := s.postTargetRepo.CreateBatch(targets); err != nil {
		return nil, errors.New("failed to create post targets")
	}

	return &PostWithTargets{Post: *post, Targets: targets}, nil
}

func (s *PostService) List(userID string, page, limit int) ([]model.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.postRepo.FindByUserID(userID, page, limit)
}

func (s *PostService) GetByID(userID, id string) (*PostWithTargets, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("post not found")
	}
	if post.UserID != userID {
		return nil, errors.New("post not found")
	}
	targets, err := s.postTargetRepo.FindByPostID(id)
	if err != nil {
		return nil, errors.New("failed to fetch post targets")
	}
	return &PostWithTargets{Post: *post, Targets: targets}, nil
}

func (s *PostService) Update(userID, id string, req UpdatePostRequest) (*model.Post, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("post not found")
	}
	if post.UserID != userID {
		return nil, errors.New("post not found")
	}
	if post.Status != model.PostStatusDraft {
		return nil, errors.New("only draft posts can be updated")
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.ScheduledAt != nil {
		updates["scheduled_at"] = req.ScheduledAt
		updates["status"] = model.PostStatusScheduled
	}

	if err := s.postRepo.Update(id, updates); err != nil {
		return nil, errors.New("failed to update post")
	}
	return s.postRepo.FindByID(id)
}

func (s *PostService) Delete(userID, id string) error {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return errors.New("post not found")
	}
	if post.UserID != userID {
		return errors.New("post not found")
	}
	return s.postRepo.SoftDelete(id)
}

func (s *PostService) Publish(userID, id string) (*PostWithTargets, error) {
	post, err := s.postRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("post not found")
	}
	if post.UserID != userID {
		return nil, errors.New("post not found")
	}
	if post.Status == model.PostStatusPublished {
		return nil, errors.New("post already published")
	}

	targets, err := s.postTargetRepo.FindByPostID(id)
	if err != nil {
		return nil, errors.New("failed to fetch post targets")
	}

	s.publishToTargets(post, targets)

	updated, _ := s.postRepo.FindByID(id)
	updatedTargets, _ := s.postTargetRepo.FindByPostID(id)
	return &PostWithTargets{Post: *updated, Targets: updatedTargets}, nil
}

func (s *PostService) PublishScheduled() {
	posts, err := s.postRepo.FindScheduledDue()
	if err != nil {
		return
	}
	for i := range posts {
		targets, err := s.postTargetRepo.FindByPostID(posts[i].ID)
		if err != nil {
			continue
		}
		s.publishToTargets(&posts[i], targets)
	}
}

func (s *PostService) publishToTargets(post *model.Post, targets []model.PostTarget) {
	allPublished := true
	for _, target := range targets {
		account, err := s.socialAccountRepo.FindByID(target.SocialAccountID)
		if err != nil {
			errMsg := "social account not found"
			_ = s.postTargetRepo.UpdateStatus(target.ID, model.PostTargetStatusFailed, &errMsg)
			allPublished = false
			continue
		}

		err = s.threadsClient.Publish(account.ThreadsUserID, account.AccessToken, post.Content)
		if err != nil {
			errMsg := err.Error()
			_ = s.postTargetRepo.UpdateStatus(target.ID, model.PostTargetStatusFailed, &errMsg)
			allPublished = false
			continue
		}

		_ = s.postTargetRepo.UpdateStatus(target.ID, model.PostTargetStatusPublished, nil)
	}

	now := time.Now()
	if allPublished {
		_ = s.postRepo.Update(post.ID, map[string]interface{}{
			"status":       model.PostStatusPublished,
			"published_at": now,
		})
	}
}
