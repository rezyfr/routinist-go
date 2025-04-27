package http

import (
	v1 "routinist/internal/controller/v1"
	"routinist/internal/middleware"
	"routinist/internal/usecase"
	"routinist/pkg/logger"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	handler *gin.Engine,
	l logger.Interface,
	tAuth usecase.Auth,
) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	h := handler.Group("/api/v1")
	h.Use(middleware.ContentTypeApplicationJson())

	{
		v1.NewAuthRoutes(h, tAuth, l)
	}
}
