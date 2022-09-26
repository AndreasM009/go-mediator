package mediator

import "reflect"

type RequestOptionInterface interface {
	GetRequestType() reflect.Type
	GetRequestHandler() RequestHandlerInterface
}

type RequestOption[TRequest IsRequest, TResponse IsResponse, THandler IsRequestHandler[TRequest, TResponse]] struct {
	genericHandler THandler
	requestType    reflect.Type
}

func WithRequest[TRequest IsRequest, TResponse IsResponse, THandler IsRequestHandler[TRequest, TResponse]](req TRequest, handler THandler) RequestOptionInterface {
	return &RequestOption[TRequest, TResponse, THandler]{
		requestType:    reflect.TypeOf(req),
		genericHandler: handler,
	}
}

func (o *RequestOption[TRequest, TResponse, THandler]) GetRequestHandler() RequestHandlerInterface {
	return o.genericHandler
}

func (o *RequestOption[TRequest, TResponse, THandler]) GetRequestType() reflect.Type {
	return o.requestType
}
