package mediator

import (
	"context"
	"fmt"
	"reflect"
)

type Mediator struct {
	requestHandlers      map[reflect.Type]RequestHandlerInterface
	behaviors            []PipelineBehavior
	notificationHandlers map[reflect.Type][]NotificationHandlerInterface
}

func NewMediator() *Mediator {
	return &Mediator{
		requestHandlers:      make(map[reflect.Type]RequestHandlerInterface),
		behaviors:            nil,
		notificationHandlers: make(map[reflect.Type][]NotificationHandlerInterface),
	}
}

func (m *Mediator) ConfigureRequests(opts ...RequestOptionInterface) {
	for _, opt := range opts {
		m.requestHandlers[opt.GetRequestType()] = opt.GetRequestHandler()
	}
}

func (m *Mediator) ConfigureBehaviors(opts ...PipelineBehaviorOption) {
	if len(opts) == 0 {
		m.behaviors = make([]PipelineBehavior, 0)
		return
	}

	m.behaviors = make([]PipelineBehavior, len(opts))

	for i, opt := range opts {
		m.behaviors[i] = opt.processor
	}
}

func (m *Mediator) ConfigureNotifications(opts ...NotificationOptionInterface) {
	for _, opt := range opts {
		if hs, ok := m.notificationHandlers[opt.GetNotificationType()]; ok {
			hs = append(hs, opt.GetNotificationHandler())
			m.notificationHandlers[opt.GetNotificationType()] = hs
		} else {
			hs = make([]NotificationHandlerInterface, 1)
			hs[0] = opt.GetNotificationHandler()
			m.notificationHandlers[opt.GetNotificationType()] = hs
		}
	}
}

func Send[TRequest IsRequest, TResponse IsResponse](ctx context.Context, m *Mediator, req TRequest) (TResponse, error) {
	count := len(m.behaviors)
	idx := 0
	var handler RequestHandler[TRequest, TResponse]

	if h, ok := m.requestHandlers[reflect.TypeOf(req)]; ok {
		handler = h.(RequestHandler[TRequest, TResponse])
	} else {
		var result TResponse
		err := fmt.Errorf("no handler found for request of type %s", reflect.TypeOf(req).String())
		return result, err
	}

	var next func(ctx context.Context, req Request) (Response, error) = nil

	next = func(ctx context.Context, req Request) (Response, error) {
		if idx >= count {
			return handler.Handle(ctx, req.(TRequest))
		}

		p := m.behaviors[idx]
		idx++
		return p.Process(ctx, req, next)
	}

	r, err := next(ctx, req)
	result := r.(TResponse)

	return result, err
}

func Publish[TNotification IsNotification](ctx context.Context, m *Mediator, notification TNotification) error {
	errResult := []string{}

	if hs, ok := m.notificationHandlers[reflect.TypeOf(notification)]; ok {
		for _, h := range hs {
			if err := h.(NotificationHandler[TNotification]).Handle(ctx, notification); err != nil {
				errResult = append(errResult, err.Error())
			}
		}

		if len(errResult) > 0 {
			return fmt.Errorf("notfication errors: %v", errResult)
		}
		return nil
	} else {
		err := fmt.Errorf("no handler found for notification of type %s", reflect.TypeOf(notification).String())
		return err
	}
}
