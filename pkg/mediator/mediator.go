package mediator

import (
	"context"
	"fmt"
	"reflect"
)

type Mediator struct {
	requestHandlers             map[reflect.Type]RequestHandler
	preRequestProcessors        []PreRequestProcessor
	postRequestProcessors       []PostRequestProcessor
	preRequestProcessorWrapper  *preRequestProcessorWrapper
	postRequestprocessorWrapper *postRequestProcessorWrapper
}

func NewMediator() *Mediator {
	return &Mediator{
		requestHandlers: make(map[reflect.Type]RequestHandler),
	}
}

func (m *Mediator) ConfigureRequests(opts ...*RequestOption) *Mediator {
	for _, opt := range opts {
		opt.configure(m)
	}

	return m
}

func (m *Mediator) ConfigureRequestProcessors(opts ...RequestProcessorOption) *Mediator {
	for _, opt := range opts {
		opt(m)
	}

	m.preRequestProcessorWrapper = buildPreRequestProcessorWrapper(m.preRequestProcessors)
	m.postRequestprocessorWrapper = buildPostRequestprocessorWrapper(m.postRequestProcessors)
	return m
}

func (m *Mediator) Send(ctx context.Context, rq Request) <-chan Response {
	r := make(chan Response)

	if h, ok := m.requestHandlers[reflect.TypeOf(rq)]; ok {
		hw := buildRequestHandlerWrapper(m.preRequestProcessorWrapper, h, m.postRequestprocessorWrapper)
		go func() {
			defer close(r)
			r <- hw.Handle(ctx, rq)
		}()

		return r
	}

	go func() {
		defer close(r)
		r <- CreateEmtpyResponse(fmt.Errorf("Error: no handler found for %s", reflect.TypeOf(rq)))
	}()

	return r
}
