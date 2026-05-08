package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qiblatdigital/zoztool-api/internal/helper"
	"github.com/qiblatdigital/zoztool-api/internal/service"
)

type SocialAccountHandler struct {
	svc *service.SocialAccountService
}

func NewSocialAccountHandler(svc *service.SocialAccountService) *SocialAccountHandler {
	return &SocialAccountHandler{svc: svc}
}

func (h *SocialAccountHandler) Create(c *gin.Context) {
	var req service.CreateSocialAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	userID := c.GetString("user_id")
	account, err := h.svc.Create(userID, req)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	helper.Success(c, http.StatusCreated, "social account created", account)
}

func (h *SocialAccountHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	accounts, err := h.svc.List(userID)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "failed to fetch social accounts")
		return
	}
	helper.Success(c, http.StatusOK, "social accounts retrieved", accounts)
}

func (h *SocialAccountHandler) GetByID(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")
	account, err := h.svc.GetByID(userID, id)
	if err != nil {
		helper.Error(c, http.StatusNotFound, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "social account retrieved", account)
}

func (h *SocialAccountHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")
	if err := h.svc.Delete(userID, id); err != nil {
		helper.Error(c, http.StatusNotFound, err.Error())
		return
	}
	helper.Success(c, http.StatusOK, "social account deleted", nil)
}
