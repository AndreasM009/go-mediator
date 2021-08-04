package mediator

import (
	"context"
	"fmt"
	"reflect"
)

type Mediator struct {
	requestHandlers             map[reflect.Type]RequestHandler
	notificationHandlers        map[reflect.Type][]NotificationHandler
	preRequestProcessors        []PreRequestProcessor
	postRequestProcessors       []PostRequestProcessor
	preRequestProcessorWrapper  *preRequestProcessorWrapper
	postRequestprocessorWrapper *postRequestProcessorWrapper
}

func NewMediator() *Mediator {
	return &Mediator{
		requestHandlers:      make(map[reflect.Type]RequestHandler),
		notificationHandlers: make(map[reflect.Type][]NotificationHandler),
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

func (m *Mediator) ConfigureNotifications(opts ...NotificationOption) *Mediator {
	for _, opt := range opts {
		opt(m)
	}

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

func (m *Mediator) Publish(ctx context.Context, n Notification) <-chan error {
	r := make(chan error)
	var result error = nil

	if hs, ok := m.notificationHandlers[reflect.TypeOf(n)]; ok {
		go func() {
			defer close(r)

			chans := make([]chan error, len(hs))

			for i, h := range hs {
				ch := make(chan error)
				chans[i] = ch
				th := h

				go func() {
					defer close(ch)
					err := th.Handle(ctx, n)
					if nil != err {
						ch <- err
					}
				}()
			}

			// wait until all notifications are processed
			for _, ch := range chans {
				err := <-ch
				if err != nil {
					result = fmt.Errorf("%s\n%s", result, err)
				}
			}
			r <- result
		}()

		return r
	}

	go func() {
		defer close(r)
		r <- fmt.Errorf("Error: no handler found for %s", reflect.TypeOf(n))
	}()

	return r
}
