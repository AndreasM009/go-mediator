package mediator

import "reflect"

type NotificationOption func(m *Mediator)

func WithNotification(n Notification, h NotificationHandler) NotificationOption {
	return func(m *Mediator) {
		t := reflect.TypeOf(n)
		if val, ok := m.notificationHandlers[t]; ok {
			m.notificationHandlers[t] = append(val, h)
			return
		}

		hs := []NotificationHandler{h}
		m.notificationHandlers[t] = hs
	}
}
