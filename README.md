# Mediator

A simple mediator implementation in GO. 

It implements a in-process messaging system to send either request/response messages or publish and dispatch messages.


## Send request/response

Send a simple message for a dedicated handler (request/response) with no return value:

```go
m := mediator.NewMediator()
ctx := context.Background()
invoked := false

type testRequest struct {
}

handler := func(ctx context.Context, rq mediator.Request) mediator.Response {
    invoked = true
    return mediator.CreateEmtpyResponse(nil)
}

m.ConfigureRequests(
    mediator.WithRequest(&testRequest{}, mediator.RequestHandlerFunc(handler)))

r := <-m.Send(ctx, &testRequest{})

if r.HasError() {
    //...
}
```

Send a simple message for a dedicated handler (request/response) with return value:

```go
m := mediator.NewMediator()
ctx := context.Background()
invoked := false

type testRequest struct {
}

handler := func(ctx context.Context, rq mediator.Request) mediator.Response {
    invoked = true
    return mediator.CreateResponse(true, nil)
}

m.ConfigureRequests(
    mediator.WithRequest(&testRequest{}, mediator.RequestHandlerFunc(handler)))

r := <-m.Send(ctx, &testRequest{})

if r.HasError() {
    // ...
}

fmt.Println(r.Result().(bool))
```

## Pre- and Post-Processors

Sending in request/response mode, the mediator allows you to add Pre- and Post-Processors to create a process pipeline:

```go
m := mediator.NewMediator()
ctx := context.Background()
handlerInvoked := false
preProcessorInvokked := false
postProcessorInvoked := false
postWasTheLast := true

type testRequest struct {
}

handler := func(ctx context.Context, rq mediator.Request) mediator.Response {
    handlerInvoked = true
    return mediator.CreateEmtpyResponse(nil)
}

preProcessor := func(ctx context.Context, rq mediator.Request, next mediator.NextPreRequestProcessorDelegate) error {
    preProcessorInvokked = true
    postWasTheLast = false
    return next(ctx, rq)
}

postProcessor := func(ctx context.Context, rq mediator.Request, resp mediator.Response, next mediator.NextPostRequestProcessorDelegate) error {
    postProcessorInvoked = true
    postWasTheLast = true
    return next(ctx, rq, resp)
}

m.
    ConfigureRequests(
        mediator.WithRequest(&testRequest{}, mediator.RequestHandlerFunc(handler))).
    ConfigureRequestProcessors(
        mediator.WithPreRequestProcessor(mediator.PreRequestProcessorFunc(preProcessor)),
        mediator.WithPostRequestProcessor(mediator.PostRequestProcessorFunc(postProcessor)))

r := <-m.Send(ctx, &testRequest{})
```

Pre- and Post-Processors can be defined for each request, too:

```go
m.
	ConfigureRequests(
		mediator.WithRequest(&testRequest{}, mediator.RequestHandlerFunc(handler)).
			WithPreProcessor(mediator.PreRequestProcessorFunc(preProcessor)).
			WithPostProcessor(mediator.PostRequestProcessorFunc(postProcessor))).
	ConfigureRequestProcessors(
		mediator.WithPreRequestProcessor(mediator.PreRequestProcessorFunc(preMediatorProcessor)),
		mediator.WithPostRequestProcessor(mediator.PostRequestProcessorFunc(postMediatorProcessor)))

r := <-m.Send(ctx, &testRequest{})
```

## Publish

Publish a notification to multiple handlers:

```go
m := mediator.NewMediator()
ctx := context.Background()

handlerOneInvoked := false
handlerTwoInvoked := false

type testNotification struct {
}

handlerOne := func(ctx context.Context, n mediator.Notification) error {
    handlerOneInvoked = true
    return nil
}

handlerTwo := func(ctx context.Context, n mediator.Notification) error {
    handlerTwoInvoked = true
    return nil
}

m.
    ConfigureNotifications(
        mediator.WithNotification(&testNotification{}, mediator.NotificationHandlerFunc(handlerOne)),
        mediator.WithNotification(&testNotification{}, mediator.NotificationHandlerFunc(handlerTwo)))

err := <-m.Publish(ctx, &testNotification{})

assert.Nil(t, err)
assert.True(t, handlerOneInvoked)
assert.True(t, handlerTwoInvoked)
```