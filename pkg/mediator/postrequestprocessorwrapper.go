package mediator

import "context"

type postRequestProcessorWrapper struct {
	processors []PostRequestProcessor
}

func emptyPostProcessor(ctx context.Context, rq Request, resp Response, next NextPostRequestProcessorDelegate) error {
	return nil
}

func buildPostRequestprocessorWrapper(ps []PostRequestProcessor) *postRequestProcessorWrapper {
	return &postRequestProcessorWrapper{
		processors: append(ps, PostRequestProcessorFunc(emptyPostProcessor)),
	}
}

func (pw *postRequestProcessorWrapper) Process(ctx context.Context, rq Request, resp Response) error {
	cnt := len(pw.processors)
	idx := 0

	var next func(ctx context.Context, rq Request, resp Response) error = nil

	next = func(ctx context.Context, rq Request, resp Response) error {
		if idx >= cnt {
			return nil
		}

		p := pw.processors[idx]
		idx++

		return p.Process(ctx, rq, resp, next)
	}

	return next(ctx, rq, resp)
}
