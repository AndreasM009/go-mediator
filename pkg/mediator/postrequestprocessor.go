package mediator

import "context"

type PostRequestProcessor interface {
	Process(ctx context.Context, rq Request, resp Response, next NextPostRequestProcessorDelegate) error
}

type NextPostRequestProcessorDelegate func(ctx context.Context, rq Request, resp Response) error

type PostRequestProcessorFunc func(ctx context.Context, rq Request, resp Response, next NextPostRequestProcessorDelegate) error

func (f PostRequestProcessorFunc) Process(ctx context.Context, rq Request, resp Response, next NextPostRequestProcessorDelegate) error {
	return f(ctx, rq, resp, next)
}
