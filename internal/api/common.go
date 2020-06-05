package api

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/duchenhao/backend-demo/internal/model"
)

// func handler(ctx *model.ReqContext) Response
// func handler(ctx *model.ReqContext, params *RequestStruct) Response
func Wrap(handler interface{}) gin.HandlerFunc {
	t := reflect.TypeOf(handler)
	if t.Kind() != reflect.Func {
		panic("handler should be function type")
	}

	fnNumIn := t.NumIn()
	if fnNumIn == 0 || fnNumIn > 2 {
		panic("handler require 1 or 2 input params")
	}

	tc := reflect.TypeOf(&model.ReqContext{})
	if t.In(0) != tc {
		panic("handler first param should by type of *model.ReqContext")
	}

	if t.NumOut() != 1 {
		panic("handler return values should contain response data")
	}

	return func(ctx *gin.Context) {
		reqCtxI, _ := ctx.Get("ctx")
		reqCtx := reqCtxI.(*model.ReqContext)

		params := make([]reflect.Value, fnNumIn)
		params[0] = reflect.ValueOf(reqCtx)
		if fnNumIn == 2 {
			req := reflect.New(t.In(1).Elem()).Interface()
			if err := ctx.ShouldBind(req); err != nil {
				reqCtx.Logger.Error(err.Error())
				ParamsError().WriteTo(reqCtx)
				return
			}
			params[1] = reflect.ValueOf(req)
			reqCtx.Logger.With(zap.Any("form", req))
		}

		ret := reflect.ValueOf(handler).Call(params)
		var res Response
		if ret[0].IsNil() {
			res = ServerError()
		} else {
			res = ret[0].Interface().(Response)
		}
		res.WriteTo(reqCtx)
	}
}

type Response interface {
	WriteTo(ctx *model.ReqContext)
}

type JsonResponse struct {
	err    error
	data   interface{}
	status int
}

func (r *JsonResponse) WriteTo(ctx *model.ReqContext) {
	ctx.JSON(r.status, r.data)
}

func Error(status int, msg string) *JsonResponse {
	data := gin.H{
		"message": msg,
	}
	return &JsonResponse{
		data:   data,
		status: status,
	}
}

func ServerError() *JsonResponse {
	return Error(500, "Internal Error")
}

func ParamsError() *JsonResponse {
	return Error(400, "Params Error")
}

func AuthError() *JsonResponse {
	return Error(401, "Auth Error")
}

func JSON(data interface{}) *JsonResponse {
	return &JsonResponse{
		err:    nil,
		data:   data,
		status: 200,
	}
}
