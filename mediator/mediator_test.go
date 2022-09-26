package mediator

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	pipelineBehaviorFuncCalled    = false
	notificationHandlerFuncCalled = false
)

func Echo(ctx context.Context, msg string) (string, error) {
	return msg, nil
}

type EchoStruct struct {
}

func (e *EchoStruct) Handle(ctx context.Context, msg string) (string, error) {
	return msg, nil
}

type EchoValueStruct struct {
}

func (e EchoValueStruct) Handle(ctx context.Context, msg string) (string, error) {
	return msg, nil
}

func BehaviorFunc(ctx context.Context, req Request, next NextBipelineBehavior) (Response, error) {
	pipelineBehaviorFuncCalled = true
	return next(ctx, req)
}

func NotificationFunc(ctx context.Context, msg string) error {
	notificationHandlerFuncCalled = true
	return nil
}

type EventStruct struct {
	Message string
}

type NotificationHandlerStruct struct{}

func (h *NotificationHandlerStruct) Handle(ctx context.Context, notification EventStruct) error {
	fmt.Println(notification.Message)
	return nil
}

func TestRequestFunc(t *testing.T) {
	m := NewMediator()

	m.ConfigureRequests(
		WithRequest[string, string]("", RequestHandlerFunc[string, string](Echo)))

	echo, _ := Send[string, string](context.TODO(), m, "Hello World")
	assert.Equal(t, "Hello World", echo)

	fmt.Println(echo)
}

func TestRequestStructHandler(t *testing.T) {
	m := NewMediator()

	m.ConfigureRequests(
		WithRequest[string, string]("", &EchoStruct{}))

	echo, _ := Send[string, string](context.TODO(), m, "Hello World")
	assert.Equal(t, "Hello World", echo)

	fmt.Println(echo)
}

func TestRequestValueStructHandler(t *testing.T) {
	m := NewMediator()

	m.ConfigureRequests(
		WithRequest[string, string]("", EchoValueStruct{}))

	echo, _ := Send[string, string](context.TODO(), m, "Hello World")
	assert.Equal(t, "Hello World", echo)

	fmt.Println(echo)
}

func TestBehaviorFunc(t *testing.T) {
	m := NewMediator()

	m.ConfigureRequests(
		WithRequest[string, string]("", EchoValueStruct{}))

	m.ConfigureBehaviors(
		WithBehavior(PipelineBehaviorFunc(BehaviorFunc)))

	pipelineBehaviorFuncCalled = false
	echo, _ := Send[string, string](context.TODO(), m, "Hello World")
	assert.Equal(t, "Hello World", echo)
	assert.Equal(t, true, pipelineBehaviorFuncCalled)

	fmt.Println(echo)
}

func TestNotificationFunc(t *testing.T) {
	m := NewMediator()

	m.ConfigureNotifications(
		WithNotification("", NotificationHandlerFunc[string](NotificationFunc)))

	notificationHandlerFuncCalled = false

	_ = Publish(context.TODO(), m, "Hello World")

	assert.Equal(t, true, notificationHandlerFuncCalled)
}

func TestNotificationStruct(t *testing.T) {
	m := NewMediator()

	m.ConfigureNotifications(
		WithNotification(EventStruct{}, &NotificationHandlerStruct{}))

	err := Publish(context.TODO(), m, EventStruct{"Hello World"})
	assert.Nil(t, err)
}
