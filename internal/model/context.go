package model

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ReqContext struct {
	*gin.Context

	RequestId    string
	SignedInUser *SignedInUser

	Logger *zap.Logger
}
