package mediator

type Publisher interface {
	Publish(n Notification) <-chan error
}
