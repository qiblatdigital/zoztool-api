package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qiblatdigital/zoztool-api/internal/helper"
	"github.com/qiblatdigital/zoztool-api/internal/model"
	"github.com/qiblatdigital/zoztool-api/internal/service"
)

const refreshTokenCookie = "refresh_token"

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.authSvc.Register(req)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusCreated, "registration successful", model.ToUserResponse(user))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, tokens, err := h.authSvc.Login(req)
	if err != nil {
		helper.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Set refresh token in httpOnly cookie (7 days)
	c.SetCookie(refreshTokenCookie, tokens.RefreshToken, 7*24*3600, "/", "", false, true)

	helper.Success(c, http.StatusOK, "login successful", gin.H{
		"access_token": tokens.AccessToken,
		"user":         model.ToUserResponse(user),
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie(refreshTokenCookie, "", -1, "/", "", false, true)
	helper.Success(c, http.StatusOK, "logout successful", nil)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshTokenCookie)
	if err != nil || refreshToken == "" {
		helper.Error(c, http.StatusUnauthorized, "refresh token not found")
		return
	}

	accessToken, err := h.authSvc.RefreshToken(refreshToken)
	if err != nil {
		helper.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "token refreshed", gin.H{
		"access_token": accessToken,
	})
}
