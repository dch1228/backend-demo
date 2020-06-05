package api

import (
	"github.com/gin-gonic/gin"

	"github.com/duchenhao/backend-demo/internal/api/forms"
	"github.com/duchenhao/backend-demo/internal/bus"
	"github.com/duchenhao/backend-demo/internal/model"
)

func signUp(ctx *model.ReqContext, form *forms.SignUpForm) Response {
	if form.Password != form.Password2 {
		ctx.Logger.Info("password not match")
		return Error(400, "密码不匹配")
	}

	existing := &model.GetUserByNameQuery{Name: form.Name}
	if err := bus.Dispatch(ctx, existing); err != nil && err != model.ErrUserNotFound {
		ctx.Logger.Error(err.Error())
		return ServerError()
	}
	if existing.User != nil {
		ctx.Logger.Info("user exists")
		return Error(400, "用户已存在")
	}

	cmd := &model.CreateUserCommand{}
	cmd.Name = form.Name
	cmd.Password = form.Password
	if err := bus.Dispatch(ctx, cmd); err != nil {
		ctx.Logger.Error(err.Error())
		return ServerError()
	}

	user := cmd.User
	tokenCmd := &model.CreateTokenCommand{
		User: user,
	}
	if err := bus.Dispatch(ctx, tokenCmd); err != nil {
		ctx.Logger.Error(err.Error())
		return ServerError()
	}

	res := gin.H{
		"access_token":  tokenCmd.AccessToken,
		"refresh_token": tokenCmd.RefreshToken,
	}
	return JSON(res)
}

func login(ctx *model.ReqContext, form *forms.LoginForm) Response {
	query := &model.LoginQuery{}
	query.Name = form.Name
	query.Password = form.Password
	if err := bus.Dispatch(ctx, query); err != nil {
		ctx.Logger.Error(err.Error())
		if err == model.ErrInvalidPassword {
			return AuthError()
		}
		return ServerError()
	}

	tokenCmd := &model.CreateTokenCommand{
		User: query.User,
	}
	if err := bus.Dispatch(ctx, tokenCmd); err != nil {
		ctx.Logger.Error(err.Error())
		return ServerError()
	}

	res := gin.H{
		"access_token":  tokenCmd.AccessToken,
		"refresh_token": tokenCmd.RefreshToken,
	}
	return JSON(res)
}

func refreshToken(ctx *model.ReqContext, form *forms.RefreshTokenForm) Response {
	query := &model.RefreshTokenCommand{}
	query.RefreshToken = form.Token
	if err := bus.Dispatch(ctx, query); err != nil {
		ctx.Logger.Error(err.Error())
		return ServerError()
	}

	res := gin.H{
		"access_token": query.AccessToken,
	}
	return JSON(res)
}
