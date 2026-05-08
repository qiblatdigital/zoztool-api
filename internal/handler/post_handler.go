package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qiblatdigital/zoztool-api/internal/helper"
	"github.com/qiblatdigital/zoztool-api/internal/service"
)

type PostHandler struct {
	svc *service.PostService
}

func NewPostHandler(svc *service.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

func (h *PostHandler) Create(c *gin.Context) {
	var req service.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	userID := c.GetString("user_id")
	post, err := h.svc.Create(userID, req)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	helper.Success(c, http.StatusCreated, "post created", post)
}

func (h *PostHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	posts, total, err := h.svc.List(userID, page, limit)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "failed to fetch posts")
		return
	}
	helper.SuccessPaginated(c, "posts retrieved", posts, page, limit, int(total))
}

func (h *PostHandler) GetByID(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")
	post, err := h.svc.GetByID(userID, id)
	if err != nil {
		helper.Error(c, http.StatusNotFound, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "post retrieved", post)
}

func (h *PostHandler) Update(c *gin.Context) {
	var req service.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	userID := c.GetString("user_id")
	id := c.Param("id")
	post, err := h.svc.Update(userID, id, req)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "post updated", post)
}

func (h *PostHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")
	if err := h.svc.Delete(userID, id); err != nil {
		helper.Error(c, http.StatusNotFound, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "post deleted", nil)
}

func (h *PostHandler) Publish(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")
	post, err := h.svc.Publish(userID, id)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "post published", post)
}
