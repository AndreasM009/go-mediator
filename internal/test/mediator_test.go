package test

import (
	"context"
	"testing"

	"github.com/andreasM009/go-mediator/pkg/mediator"
	"github.com/stretchr/testify/assert"
)

func TestRequestWithNoResult(t *testing.T) {
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

	assert.False(t, r.HasError())
	assert.Nil(t, r.Error())
	assert.Nil(t, r.Result())
	assert.True(t, invoked)
}

func TestRequestWithResult(t *testing.T) {
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

	assert.False(t, r.HasError())
	assert.Nil(t, r.Error())
	assert.True(t, r.Result().(bool))
	assert.True(t, invoked)
}

func TestMediatorRequestPreProcessor(t *testing.T) {
	m := mediator.NewMediator()
	ctx := context.Background()
	handlerInvoked := false
	preProcessorInvokked := false

	type testRequest struct {
	}

	handler := func(ctx context.Context, rq mediator.Request) mediator.Response {
		handlerInvoked = true
		return mediator.CreateEmtpyResponse(nil)
	}

	preProcessor := func(ctx context.Context, rq mediator.Request, next mediator.NextPreRequestProcessorDelegate) error {
		preProcessorInvokked = true
		return next(ctx, rq)
	}

	m.
		ConfigureRequests(
			mediator.WithRequest(&testRequest{}, mediator.RequestHandlerFunc(handler))).
		ConfigureRequestProcessors(
			mediator.WithPreRequestProcessor(mediator.PreRequestProcessorFunc(preProcessor)))

	r := <-m.Send(ctx, &testRequest{})

	assert.False(t, r.HasError())
	assert.Nil(t, r.Error())
	assert.Nil(t, r.Result())
	assert.True(t, handlerInvoked)
	assert.True(t, preProcessorInvokked)
}

func TestMediatorRequestPostprocessor(t *testing.T) {
	m := mediator.NewMediator()
	ctx := context.Background()
	handlerInvoked := false
	postProcessorInvoked := false

	type testRequest struct {
	}

	handler := func(ctx context.Context, rq mediator.Request) mediator.Response {
		handlerInvoked = true
		return mediator.CreateEmtpyResponse(nil)
	}

	postProcessor := func(ctx context.Context, rq mediator.Request, resp mediator.Response, next mediator.NextPostRequestProcessorDelegate) error {
		postProcessorInvoked = true
		return next(ctx, rq, resp)
	}

	m.
		ConfigureRequests(
			mediator.WithRequest(&testRequest{}, mediator.RequestHandlerFunc(handler))).
		ConfigureRequestProcessors(
			mediator.WithPostRequestProcessor(mediator.PostRequestProcessorFunc(postProcessor)))

	r := <-m.Send(ctx, &testRequest{})

	assert.False(t, r.HasError())
	assert.Nil(t, r.Error())
	assert.Nil(t, r.Result())
	assert.True(t, handlerInvoked)
	assert.True(t, postProcessorInvoked)
}

func TestMediatorRequestProcessor(t *testing.T) {
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

	assert.False(t, r.HasError())
	assert.Nil(t, r.Error())
	assert.Nil(t, r.Result())
	assert.True(t, handlerInvoked)
	assert.True(t, preProcessorInvokked)
	assert.True(t, postProcessorInvoked)
	assert.True(t, postWasTheLast)
}

func TestMediatorAndRequestProcessors(t *testing.T) {
	m := mediator.NewMediator()
	ctx := context.Background()
	handlerInvoked := false

	preMediatorProcessorInvokked := false
	postMediatorProcessorInvoked := false
	postMediatorWasTheLast := true

	preProcessorInvokked := false
	postProcessorInvoked := false
	postWasTheLast := true

	lastMustBePostMediatorProcessor := false

	type testRequest struct {
	}

	handler := func(ctx context.Context, rq mediator.Request) mediator.Response {
		handlerInvoked = true
		return mediator.CreateEmtpyResponse(nil)
	}

	preMediatorProcessor := func(ctx context.Context, rq mediator.Request, next mediator.NextPreRequestProcessorDelegate) error {
		preMediatorProcessorInvokked = true
		postMediatorWasTheLast = false
		lastMustBePostMediatorProcessor = false
		return next(ctx, rq)
	}

	postMediatorProcessor := func(ctx context.Context, rq mediator.Request, resp mediator.Response, next mediator.NextPostRequestProcessorDelegate) error {
		postMediatorProcessorInvoked = true
		postMediatorWasTheLast = true
		lastMustBePostMediatorProcessor = true
		return next(ctx, rq, resp)
	}

	preProcessor := func(ctx context.Context, rq mediator.Request, next mediator.NextPreRequestProcessorDelegate) error {
		preProcessorInvokked = true
		postWasTheLast = false
		lastMustBePostMediatorProcessor = false
		return next(ctx, rq)
	}

	postProcessor := func(ctx context.Context, rq mediator.Request, resp mediator.Response, next mediator.NextPostRequestProcessorDelegate) error {
		postProcessorInvoked = true
		postWasTheLast = true
		lastMustBePostMediatorProcessor = false
		return next(ctx, rq, resp)
	}

	m.
		ConfigureRequests(
			mediator.WithRequest(&testRequest{}, mediator.RequestHandlerFunc(handler)).
				WithPreProcessor(mediator.PreRequestProcessorFunc(preProcessor)).
				WithPostProcessor(mediator.PostRequestProcessorFunc(postProcessor))).
		ConfigureRequestProcessors(
			mediator.WithPreRequestProcessor(mediator.PreRequestProcessorFunc(preMediatorProcessor)),
			mediator.WithPostRequestProcessor(mediator.PostRequestProcessorFunc(postMediatorProcessor)))

	r := <-m.Send(ctx, &testRequest{})

	assert.False(t, r.HasError())
	assert.Nil(t, r.Error())
	assert.Nil(t, r.Result())
	assert.True(t, handlerInvoked)

	assert.True(t, preMediatorProcessorInvokked)
	assert.True(t, postMediatorProcessorInvoked)
	assert.True(t, postMediatorWasTheLast)

	assert.True(t, preProcessorInvokked)
	assert.True(t, postProcessorInvoked)
	assert.True(t, postWasTheLast)

	assert.True(t, lastMustBePostMediatorProcessor)
}
