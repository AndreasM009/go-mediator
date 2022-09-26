package mediator

import "context"

type Notification interface{}

type IsNotification interface {
	Notification
}

type NotificationHandlerInterface interface{}

type NotificationHandler[TNotification IsNotification] interface {
	NotificationHandlerInterface
	Handle(ctx context.Context, notification TNotification) error
}

type IsNotificationHandler[TNotification IsNotification] interface {
	NotificationHandler[TNotification]
}

type NotificationHandlerFunc[TNotification IsNotification] func(ctx context.Context, notification TNotification) error

func (f NotificationHandlerFunc[TNotification]) Handle(ctx context.Context, notification TNotification) error {
	return f(ctx, notification)
}
