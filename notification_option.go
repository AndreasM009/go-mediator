package mediator

import "reflect"

type NotificationOptionInterface interface {
	GetNotificationType() reflect.Type
	GetNotificationHandler() NotificationHandlerInterface
}

type NotificationOption[TNotification IsNotification, THandler IsNotificationHandler[TNotification]] struct {
	genericHandler   THandler
	notificationType reflect.Type
}

func WithNotification[TNotification IsNotification, THandler IsNotificationHandler[TNotification]](notification TNotification, handler THandler) NotificationOptionInterface {
	return &NotificationOption[TNotification, THandler]{
		genericHandler:   handler,
		notificationType: reflect.TypeOf(notification),
	}
}

func (o *NotificationOption[TNotification, THandler]) GetNotificationType() reflect.Type {
	return o.notificationType
}

func (o *NotificationOption[TNotification, THandler]) GetNotificationHandler() NotificationHandlerInterface {
	return o.genericHandler
}
