package bus

import (
	"errors"
	"reflect"

	"go.uber.org/zap"

	"github.com/duchenhao/backend-demo/internal/log"
	"github.com/duchenhao/backend-demo/internal/util"
)

type handlerFunc interface{}
type Msg interface{}

var (
	handlers  = make(map[string]handlerFunc)
	listeners = make(map[string][]handlerFunc)
	wp        = util.NewWaitGroupPool(8)
)

var ErrHandlerNotFound = errors.New("handler not found")

func AddHandler(handler handlerFunc) {
	handlerType := reflect.TypeOf(handler)
	queryTypeName := handlerType.In(0).Elem().Name()
	handlers[queryTypeName] = handler
}

func Dispatch(msg Msg) error {
	logger := log.Named("bus.Dispatch")

	msgName := reflect.TypeOf(msg).Elem().Name()
	handler := handlers[msgName]

	if handler == nil {
		logger.Error("handler not found", zap.String("msg_name", msgName))
		return ErrHandlerNotFound
	}

	params := make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(msg)

	ret := reflect.ValueOf(handler).Call(params)
	err := ret[0].Interface()
	if err == nil {
		return nil
	}
	return err.(error)
}

func AddListener(handler handlerFunc) {
	handlerType := reflect.TypeOf(handler)
	eventName := handlerType.In(1).Elem().Name()
	_, exists := listeners[eventName]
	if !exists {
		listeners[eventName] = make([]handlerFunc, 0)
	}
	listeners[eventName] = append(listeners[eventName], handler)
}

func Publish(msg Msg) {
	msgName := reflect.TypeOf(msg).Elem().Name()
	listeners := listeners[msgName]

	params := make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(msg)

	for _, listenerHandler := range listeners {
		wp.Add()
		go func(handler handlerFunc) {
			defer wp.Done()
			reflect.ValueOf(handler).Call(params)
		}(listenerHandler)
	}
}

func Close() {
	wp.Wait()
}
