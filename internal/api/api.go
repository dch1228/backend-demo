package api

import (
	"github.com/gin-gonic/gin"

	"github.com/duchenhao/backend-demo/internal/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	requireLogin := middleware.RequireLogin

	apiRoute := r.Group("/api")

	userRoute := apiRoute.Group("/user")
	{
		userRoute.POST("/sign_up", Wrap(signUp))
		userRoute.POST("/login", Wrap(login))
		userRoute.POST("/refresh_token", Wrap(refreshToken))
		userRoute.GET("/info", requireLogin, Wrap(getUserInfo))
	}
}
