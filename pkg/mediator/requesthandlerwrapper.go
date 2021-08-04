package mediator

import "context"

type requestHandlerWrapper struct {
	prewrapper  *preRequestProcessorWrapper
	postwrapper *postRequestProcessorWrapper
	handler     RequestHandler
}

func buildRequestHandlerWrapper(prewrapper *preRequestProcessorWrapper, handler RequestHandler, postwrapper *postRequestProcessorWrapper) *requestHandlerWrapper {
	return &requestHandlerWrapper{
		prewrapper:  prewrapper,
		postwrapper: postwrapper,
		handler:     handler,
	}
}

func (h *requestHandlerWrapper) Handle(ctx context.Context, rq Request) Response {
	// first run all preprocessors
	err := h.prewrapper.Process(ctx, rq)
	if err != nil {
		return CreateResponse(nil, err)
	}

	// run handler
	resp := h.handler.Handle(ctx, rq)
	if resp.HasError() {
		return resp
	}

	// run posprocessors
	err = h.postwrapper.Process(ctx, rq, resp)
	if err != nil {
		return CreateResponse(resp.Result(), err)
	}

	return resp
}
