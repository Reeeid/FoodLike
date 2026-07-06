package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"foodlike-backend/internal/adapter/handler"
)

func New(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.GET("/health", handler.NewHealthHandler(db))
	return r
}
