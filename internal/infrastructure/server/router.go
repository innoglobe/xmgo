package server

import (
	"github.com/gin-gonic/gin"
	"github.com/innoglobe/xmgo/internal/infrastructure/server/handler"
)

type RouterInterface interface {
	RegisterRoutes(secretKey string) *gin.Engine
}

type Router struct {
	companyHandler *handler.CompanyHandler
	authHandler    *handler.AuthHandler
}

func NewRouter(companyHandler *handler.CompanyHandler, authHandler *handler.AuthHandler) *Router {
	return &Router{
		companyHandler: companyHandler,
		authHandler:    authHandler,
	}
}

func (r *Router) RegisterRoutes(secretKey string) *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	r.companyHandler.RegisterRoutes(api, secretKey)

	authRoutes := router.Group("/auth")
	// For the shake of simplicity put the route inline
	authRoutes.POST("/signin", r.authHandler.SignIn)

	return router
}
