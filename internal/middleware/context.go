package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/duchenhao/backend-demo/internal/bus"
	"github.com/duchenhao/backend-demo/internal/log"
	"github.com/duchenhao/backend-demo/internal/model"
)

func ContextHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := xid.New().String()
		logger := log.Named("context").With(zap.String("request_id", requestId))

		ctx := &model.ReqContext{
			Context:   c,
			RequestId: requestId,
			Logger:    logger,
		}

		initContextWithToken(ctx)

		c.Set("ctx", ctx)
	}
}

func initContextWithToken(ctx *model.ReqContext) bool {
	token := ctx.Request.Header.Get("X-Api-Token")
	if token == "" {
		return false
	}

	cmd := &model.LookupTokenCommand{
		Token: token,
	}
	if err := bus.Dispatch(cmd); err != nil {
		ctx.Logger.Error("Failed to lookup user based on header", zap.Error(err))
		switch err.(type) {
		case jwt.ValidationError:
			jErr := err.(jwt.ValidationError)
			if jErr.Errors&jwt.ValidationErrorExpired != 0 {
				ctx.Header("WWW-Authenticate", "The access token expired")
				ctx.AbortWithStatusJSON(401, gin.H{
					"message": "The access token expired",
				})
				return false
			}
		}
		ctx.AbortWithStatusJSON(401, gin.H{
			"message": "The access token invalid",
		})
		return false
	}

	claims := cmd.Claims

	query := model.GetSignedInUserQuery{
		Ctx:    ctx,
		UserId: claims.UserId,
	}
	if err := bus.Dispatch(&query); err != nil {
		ctx.Logger.Error("Failed to get user with id", zap.String("user_id", claims.UserId), zap.Error(err))
		if err == model.ErrUserNotFound {
			ctx.AbortWithStatusJSON(400, gin.H{
				"message": err.Error(),
			})
		} else {
			ctx.AbortWithStatusJSON(500, gin.H{
				"message": "Internal Error",
			})
		}
		return false
	}

	user := query.User

	ctx.SignedInUser = user
	ctx.Logger = ctx.Logger.With(zap.String("user_id", user.UserId))

	// 30分钟更新一次
	if user.ShouldUpdateLastSeenAt() {
		ctx.Logger.Info("updating user last seen at")
		if err := bus.Dispatch(&model.UpdateUserLastSeenAtCommand{Ctx: ctx, UserId: user.UserId}); err != nil {
			ctx.Logger.Error("failed update user last seen at", zap.Error(err))
		}
	}

	return true
}
