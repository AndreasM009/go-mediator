package mediator

import "context"

type preRequestProcessorWrapper struct {
	processors []PreRequestProcessor
}

func emptyPreprocessor(ctx context.Context, rq Request, next NextPreRequestProcessorDelegate) error {
	return nil
}

func buildPreRequestProcessorWrapper(ps []PreRequestProcessor) *preRequestProcessorWrapper {
	return &preRequestProcessorWrapper{
		processors: append(ps, PreRequestProcessorFunc(emptyPreprocessor)),
	}
}

func (pw *preRequestProcessorWrapper) Process(ctx context.Context, rq Request) error {
	cnt := len(pw.processors)
	idx := 0

	var next func(ctx context.Context, rq Request) error = nil

	next = func(ctx context.Context, rq Request) error {
		if idx >= cnt {
			return nil
		}

		p := pw.processors[idx]
		idx++

		return p.Process(ctx, rq, next)
	}

	return next(ctx, rq)
}
