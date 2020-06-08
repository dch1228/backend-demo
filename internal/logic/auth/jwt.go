package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/duchenhao/backend-demo/internal/bus"
	"github.com/duchenhao/backend-demo/internal/conf"
	"github.com/duchenhao/backend-demo/internal/log"
	"github.com/duchenhao/backend-demo/internal/model"
)

var (
	logger = log.Named("jwt auth")
)

func init() {
	bus.AddHandler(createToken)
	bus.AddHandler(lookupToken)
	bus.AddHandler(refreshToken)
}

func createAccessToken(userId string) (string, error) {
	now := time.Now()
	userClaims := &model.UserClaims{
		UserId: userId,
		MapClaims: jwt.MapClaims{
			"exp": now.Add(10 * time.Minute).Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	signedAccessToken, err := accessToken.SignedString([]byte(conf.Core.Secret))
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	return signedAccessToken, nil
}

func createToken(cmd *model.CreateTokenCommand) error {
	signedAccessToken, err := createAccessToken(cmd.User.Id)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	now := time.Now()
	refreshClaims := jwt.MapClaims{
		"userId": cmd.User.Id,
		"exp":    now.Add(24 * time.Hour).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(conf.Core.Secret))
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	cmd.AccessToken = signedAccessToken
	cmd.RefreshToken = signedRefreshToken
	return nil
}

func lookupToken(cmd *model.LookupTokenCommand) error {
	token, err := jwt.ParseWithClaims(cmd.Token, &model.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.Core.Secret), nil
	})
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	if claims, ok := token.Claims.(*model.UserClaims); ok && token.Valid {
		cmd.Claims = claims
		return nil
	}

	return model.ErrTokenInvalid
}

func refreshToken(cmd *model.RefreshTokenCommand) error {
	token, err := jwt.ParseWithClaims(cmd.RefreshToken, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.Core.Secret), nil
	})
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return model.ErrTokenInvalid
	}

	userId, ok := claims["userId"].(string)
	if !ok {
		return model.ErrTokenInvalid
	}

	newToken, err := createAccessToken(userId)
	cmd.AccessToken = newToken

	return err
}
