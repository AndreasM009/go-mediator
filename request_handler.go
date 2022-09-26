package mediator

import "context"

// marker interfaces
type Request interface{}
type Response interface{}

// constraint for Request type
type IsRequest interface {
	Request
}

// constraint for Response type
type IsResponse interface {
	Response
}

// Marker interface for a request handler
type RequestHandlerInterface interface{}

// concrete request handler interface
type RequestHandler[TRequest IsRequest, TResponse IsResponse] interface {
	RequestHandlerInterface
	Handle(ctx context.Context, req TRequest) (TResponse, error)
}

// constraint for a request handler type
type IsRequestHandler[TRequest IsRequest, TResponse IsResponse] interface {
	RequestHandler[TRequest, TResponse]
}

// type for a request handler function
type RequestHandlerFunc[TRequest IsRequest, TResponse IsResponse] func(ctx context.Context, req TRequest) (TResponse, error)

// request handler function
func (f RequestHandlerFunc[TRequest, TResponse]) Handle(ctx context.Context, req TRequest) (TResponse, error) {
	return f(ctx, req)
}
