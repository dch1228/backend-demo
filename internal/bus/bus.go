package bus

import (
	"context"
	"errors"
	"reflect"
)

type handlerFunc interface{}
type Msg interface{}

var (
	handlers = make(map[string]handlerFunc)
)

var ErrHandlerNotFound = errors.New("handler not found")

func AddHandler(handler handlerFunc) {
	handlerType := reflect.TypeOf(handler)
	queryTypeName := handlerType.In(1).Elem().Name()
	handlers[queryTypeName] = handler
}

func Dispatch(ctx context.Context, msg Msg) error {
	msgName := reflect.TypeOf(msg).Elem().Name()
	handler := handlers[msgName]

	if handler == nil {
		return ErrHandlerNotFound
	}

	var params []reflect.Value
	params = append(params, reflect.ValueOf(ctx))
	params = append(params, reflect.ValueOf(msg))

	ret := reflect.ValueOf(handler).Call(params)
	err := ret[0].Interface()
	if err == nil {
		return nil
	}
	return err.(error)
}
