package users

import (
	"medisuite-api/api/handler"
	"medisuite-api/common/middlewares"

	"github.com/gin-gonic/gin"
)

type IUserRoutes interface {
	Run()
}

type UserRoutes struct {
	h handler.IHandler
	g *gin.RouterGroup
}

func NewUserRoutes(handlers handler.IHandler, group *gin.RouterGroup) *UserRoutes {
	return &UserRoutes{
		h: handlers,
		g: group,
	}
}

func (r *UserRoutes) Run() {
	groups := r.g.Group("/auth")
	{
		// routes
		groups.POST("/register", r.h.UserHandler().Register)
		groups.POST("/verify-account", r.h.UserHandler().VerifyAccount)
		groups.POST("/resend-verify", r.h.UserHandler().ResendVerify)
		groups.POST("/login", r.h.UserHandler().Login)
		groups.GET("/getuser", middlewares.AuthMiddleware(), r.h.UserHandler().GetUser)
		groups.POST("/logout", middlewares.AuthMiddleware(), r.h.UserHandler().Logout)
		groups.POST("/refresh-token", middlewares.AuthMiddleware(), r.h.UserHandler().RefreshToken)
		groups.POST("/forgot-password", r.h.UserHandler().ForgotPassword)
		groups.POST("/reset-password", r.h.UserHandler().ResetPassword)
	}
}
