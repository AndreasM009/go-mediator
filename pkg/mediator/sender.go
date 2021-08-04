package mediator

import "context"

type Sender interface {
	Send(ctx context.Context, rq Request) <-chan Response
}
