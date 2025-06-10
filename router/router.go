package router

import (
	"github.com/VaheMuradyan/Sport/user"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *user.UserHandler) {
	router := r.Group("/api")
	{
		router.POST("/register", handler.RegisterUser)
	}
}
