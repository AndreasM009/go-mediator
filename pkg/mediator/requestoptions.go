package mediator

import "reflect"

//type RequestOption func(*Mediator)
type RequestOption struct {
	requestType    reflect.Type
	handler        RequestHandler
	preprocessors  []PreRequestProcessor
	postprocessors []PostRequestProcessor
}

func WithRequest(rq Request, h RequestHandler) *RequestOption {
	return &RequestOption{
		requestType: reflect.TypeOf(rq),
		handler:     h,
	}
}

func (ro *RequestOption) configure(m *Mediator) {
	prewrapper := buildPreRequestProcessorWrapper(ro.preprocessors)

	postwrapper := buildPostRequestprocessorWrapper(ro.postprocessors)

	h := buildRequestHandlerWrapper(prewrapper, ro.handler, postwrapper)

	m.requestHandlers[ro.requestType] = h
}

func (ro *RequestOption) WithPreProcessor(p PreRequestProcessor) *RequestOption {
	ro.preprocessors = append(ro.preprocessors, p)
	return ro
}

func (ro *RequestOption) WithPostProcessor(p PostRequestProcessor) *RequestOption {
	ro.postprocessors = append(ro.postprocessors, p)
	return ro
}
