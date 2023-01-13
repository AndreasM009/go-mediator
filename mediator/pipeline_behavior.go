package mediator

import "context"

type NextPipelineBehavior func(ctx context.Context, req Request) (Response, error)

type PipelineBehavior interface {
	Process(ctx context.Context, req Request, next NextPipelineBehavior) (Response, error)
}

type PipelineBehaviorFunc func(ctx context.Context, req Request, next NextPipelineBehavior) (Response, error)

func (f PipelineBehaviorFunc) Process(ctx context.Context, req Request, next NextPipelineBehavior) (Response, error) {
	return f(ctx, req, next)
}
