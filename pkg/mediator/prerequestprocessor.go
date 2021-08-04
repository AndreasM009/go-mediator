package mediator

import "context"

type PreRequestProcessor interface {
	Process(ctx context.Context, rq Request, next NextPreRequestProcessorDelegate) error
}

type NextPreRequestProcessorDelegate func(ctx context.Context, rq Request) error

type PreRequestProcessorFunc func(ctx context.Context, rq Request, next NextPreRequestProcessorDelegate) error

func (f PreRequestProcessorFunc) Process(ctx context.Context, rq Request, next NextPreRequestProcessorDelegate) error {
	return f(ctx, rq, next)
}
