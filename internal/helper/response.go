package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, gin.H{
		"meta": gin.H{"code": code, "status": "success", "message": message},
		"data": data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"meta": gin.H{"code": code, "status": "error", "message": message},
		"data": nil,
	})
}

func SuccessPaginated(c *gin.Context, message string, data interface{}, page, limit, total int) {
	c.JSON(http.StatusOK, gin.H{
		"meta":       gin.H{"code": 200, "status": "success", "message": message},
		"data":       data,
		"pagination": gin.H{"page": page, "limit": limit, "total": total},
	})
}
