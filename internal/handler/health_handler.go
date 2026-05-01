package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qiblatdigital/zoztool-api/internal/helper"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(c *gin.Context) {
	helper.Success(c, http.StatusOK, "ok", gin.H{"status": "healthy"})
}
