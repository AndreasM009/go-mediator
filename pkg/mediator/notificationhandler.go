package mediator

import "context"

type NotificationHandler interface {
	Handle(ctx context.Context, n Notification) error
}

type NotificationHandlerFunc func(ctx context.Context, n Notification) error

func (f NotificationHandlerFunc) Handle(ctx context.Context, n Notification) error {
	return f(ctx, n)
}
