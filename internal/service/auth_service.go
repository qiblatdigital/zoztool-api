package service

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/qiblatdigital/zoztool-api/internal/model"
	"github.com/qiblatdigital/zoztool-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo      *repository.UserRepository
	jwtSecret     string
	refreshSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret, refreshSecret string) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		jwtSecret:     jwtSecret,
		refreshSecret: refreshSecret,
	}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"-"`
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (s *AuthService) Register(req RegisterRequest) (*model.User, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	existing, err := s.userRepo.FindByEmail(email)
	if err == nil && existing != nil {
		return nil, errors.New("email already registered")
	}
	// Only proceed if err is gorm.ErrRecordNotFound (user doesn't exist).
	// Any other DB error should be surfaced.
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, errors.New("failed to check existing user")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	user := &model.User{
		Email:        email,
		PasswordHash: string(hash),
		Name:         strings.TrimSpace(req.Name),
		Role:         "user",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

func (s *AuthService) Login(req LoginRequest) (*model.User, *AuthTokens, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, nil, errors.New("invalid email or password")
	}

	accessToken, err := generateToken(user, s.jwtSecret, 15*time.Minute)
	if err != nil {
		return nil, nil, errors.New("failed to generate access token")
	}

	refreshToken, err := generateToken(user, s.refreshSecret, 7*24*time.Hour)
	if err != nil {
		return nil, nil, errors.New("failed to generate refresh token")
	}

	return user, &AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	claims, err := ParseJWT(refreshToken, s.refreshSecret)
	if err != nil {
		return "", errors.New("invalid or expired refresh token")
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return "", errors.New("user not found")
	}

	accessToken, err := generateToken(user, s.jwtSecret, 15*time.Minute)
	if err != nil {
		return "", errors.New("failed to generate access token")
	}

	return accessToken, nil
}

func generateToken(user *model.User, secret string, duration time.Duration) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "zoztool-api",
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseJWT(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
