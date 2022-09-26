package mediator

import "context"

type NextBipelineBehavior func(ctx context.Context, req Request) (Response, error)

type PipelineBehavior interface {
	Process(ctx context.Context, req Request, next NextBipelineBehavior) (Response, error)
}

type PipelineBehaviorFunc func(ctx context.Context, req Request, next NextBipelineBehavior) (Response, error)

func (f PipelineBehaviorFunc) Process(ctx context.Context, req Request, next NextBipelineBehavior) (Response, error) {
	return f(ctx, req, next)
}
