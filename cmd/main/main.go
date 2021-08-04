package main

import (
	"context"
	"fmt"

	"github.com/andreasM009/go-mediator/pkg/mediator"
)

// Simple SayhelloCommand
type SayHelloCommand struct {
}

// Command with no handler registration
type NoHandlerCommand struct {
}

// Command with parameters
type EchoStringCommand struct {
	Text string
}

func NewEchoStringCommand(text string) *EchoStringCommand {
	return &EchoStringCommand{
		Text: text,
	}
}

// Command with a result of type bool
type ReturnBoolCommand struct {
}

// simple say hello notification
type SayHelloNotification struct {
}

// handle SayHelloCommand
func HandleSayHello(ctx context.Context, rq mediator.Request) mediator.Response {
	fmt.Println("Hello World!")
	return mediator.CreateEmtpyResponse(nil)
}

// handle EchoStringCommand
func HandleEchoStringCommand(ctx context.Context, rq mediator.Request) mediator.Response {
	text := rq.(*EchoStringCommand).Text
	fmt.Println(text)
	return mediator.CreateEmtpyResponse(nil)
}

// handle return bool command
func HandleReturnBoolCommand(ctx context.Context, rq mediator.Request) mediator.Response {
	fmt.Println("Handle bool command")
	return mediator.CreateResponse(true, nil)
}

func PreProcessor1(ctx context.Context, rq mediator.Request, next mediator.NextPreRequestProcessorDelegate) error {
	fmt.Println("PreProcessor1")
	return next(ctx, rq)
}

func PreProcessor2(ctx context.Context, rq mediator.Request, next mediator.NextPreRequestProcessorDelegate) error {
	fmt.Println("PreProcessor2")
	return next(ctx, rq)
}

func PostProcessor1(ctx context.Context, rq mediator.Request, resp mediator.Response, next mediator.NextPostRequestProcessorDelegate) error {
	fmt.Println("PostProcessor1")
	return next(ctx, rq, resp)
}

func PostProcessor2(ctx context.Context, rq mediator.Request, resp mediator.Response, next mediator.NextPostRequestProcessorDelegate) error {
	fmt.Println("PostProcessor2")
	return next(ctx, rq, resp)
}

func PreProcessorForAllRequests(ctx context.Context, rq mediator.Request, next mediator.NextPreRequestProcessorDelegate) error {
	fmt.Println("PreProcessorForAllRequests")
	return next(ctx, rq)
}

func PostProcessorForAllRequests(ctx context.Context, rq mediator.Request, resp mediator.Response, next mediator.NextPostRequestProcessorDelegate) error {
	fmt.Println("PostProcessorForAllRequests")
	return next(ctx, rq, resp)
}

func HandleSayHelloNotification1(ctx context.Context, n mediator.Notification) error {
	fmt.Println("Notification1: Hello World!")
	return nil
}

func HandleSayHelloNotification2(ctx context.Context, n mediator.Notification) error {
	fmt.Println("Notification2: Hello World!")
	return nil
}

func main() {
	ctx := context.Background()

	// configure mediator
	m := mediator.NewMediator().
		ConfigureRequests(
			mediator.WithRequest(&SayHelloCommand{}, mediator.RequestHandlerFunc(HandleSayHello)).
				WithPreProcessor(mediator.PreRequestProcessorFunc(PreProcessor1)).
				WithPreProcessor(mediator.PreRequestProcessorFunc(PreProcessor2)).
				WithPostProcessor(mediator.PostRequestProcessorFunc(PostProcessor1)).
				WithPostProcessor(mediator.PostRequestProcessorFunc(PostProcessor2)),
			mediator.WithRequest(&EchoStringCommand{}, mediator.RequestHandlerFunc(HandleEchoStringCommand)),
			mediator.WithRequest(&ReturnBoolCommand{}, mediator.RequestHandlerFunc(HandleReturnBoolCommand))).
		ConfigureRequestProcessors(
			mediator.WithPreRequestProcessor(mediator.PreRequestProcessorFunc(PreProcessorForAllRequests)),
			mediator.WithPostRequestProcessor(mediator.PostRequestProcessorFunc(PostProcessorForAllRequests))).
		ConfigureNotifications(
			mediator.WithNotification(&SayHelloNotification{}, mediator.NotificationHandlerFunc(HandleSayHelloNotification1)),
			mediator.WithNotification(&SayHelloNotification{}, mediator.NotificationHandlerFunc(HandleSayHelloNotification2)))

	// send
	r := <-m.Send(ctx, &SayHelloCommand{})
	if r.HasError() {
		fmt.Println(r.Error())
	}

	// send
	r = <-m.Send(ctx, &NoHandlerCommand{})
	if r.HasError() {
		fmt.Println(r.Error())
	}

	// send
	r = <-m.Send(ctx, NewEchoStringCommand("Echo: hello World!"))
	if r.HasError() {
		fmt.Println(r.Error())
	}

	// send
	r = <-m.Send(ctx, &ReturnBoolCommand{})
	if r.HasError() {
		fmt.Println(r.Error())
	}

	result := r.Result().(bool)
	fmt.Println(result)

	nr := <-m.Publish(ctx, &SayHelloNotification{})
	if nil != nr {
		fmt.Println(nr)
	}

	fmt.Println("End...")
}
