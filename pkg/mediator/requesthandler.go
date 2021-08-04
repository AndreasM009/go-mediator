package mediator

import "context"

type RequestHandler interface {
	Handle(ctx context.Context, rq Request) Response
}

type RequestHandlerFunc func(ctx context.Context, rq Request) Response

func (f RequestHandlerFunc) Handle(ctx context.Context, rq Request) Response {
	return f(ctx, rq)
}
