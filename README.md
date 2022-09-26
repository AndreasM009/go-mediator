# Mediator

A simple mediator implementation in Go.

## Send request/response 

Send a request for a dedicated handler:

```go
type EchoCommand struct {
	Message string
}

type EchoCommandHandler struct{}

func (h *EchoCommandHandler) Handle(ctx context.Context, command EchoCommand) (string, error) {
	if command.Message == "" {
		return "", errors.New("empty string provided")
	}

	return command.Message, nil
}

func main() {
	m := mediator.NewMediator()

	m.ConfigureRequests(
		mediator.WithRequest[EchoCommand, string](EchoCommand{}, &EchoCommandHandler{}))

	echo, _ := mediator.Send[EchoCommand, string](context.TODO(), m, EchoCommand{Message: "Hello World"})
	
    fmt.Println(echo)
}
```

## Pipeline behaviors

```go
type EchoCommand struct {
	Message string
}

type EchoCommandHandler struct{}

func (h *EchoCommandHandler) Handle(ctx context.Context, command EchoCommand) (string, error) {
	if command.Message == "" {
		return "", errors.New("empty string provided")
	}

	return command.Message, nil
}

func BehaviorFunc(ctx context.Context, req Request, next mediator.NextBipelineBehavior) (Response, error) {
	return next(ctx, req)
}

func main() {
	m := mediator.NewMediator()

	m.ConfigureRequests(
		mediator.WithRequest[EchoCommand, string](EchoCommand{}, &EchoCommandHandler{}))

    m.ConfigureBehaviors(
		mediator.WithBehavior(mediator.PipelineBehaviorFunc(BehaviorFunc)))

	echo, _ := Send[EchoCommand, string](context.TODO(), m, EchoCommand{Message: "Hello World"})
	
    fmt.Println(echo)
}
```

## Publish notifications

```go
type EventStruct struct {
	Message string
}

type NotificationHandlerStruct struct{}

func (h *NotificationHandlerStruct) Handle(ctx context.Context, notification EventStruct) error {
	fmt.Println(notification.Message)
	return nil
}

func main() {
	m := mediator.NewMediator()

	m.ConfigureNotifications(
		mediator.WithNotification(EventStruct{}, &NotificationHandlerStruct{}))

	_ := mediator.Publish(context.TODO(), m, EventStruct{"Hello World"})
}
```