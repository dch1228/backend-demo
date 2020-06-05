package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	apiRoute := r.Group("/api")

	userRoute := apiRoute.Group("/user")
	{
		userRoute.POST("/sign_up", Wrap(signUp))
		userRoute.POST("/login", Wrap(login))
		userRoute.POST("/refresh_token", Wrap(refreshToken))
	}
}
