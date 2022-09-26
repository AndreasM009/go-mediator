package mediator

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestHandleCommand(t *testing.T) {
	m := NewMediator()

	m.ConfigureRequests(
		WithRequest[EchoCommand, string](EchoCommand{}, &EchoCommandHandler{}))

	echo, err := Send[EchoCommand, string](context.TODO(), m, EchoCommand{Message: "Hello World"})
	assert.Nil(t, err)
	assert.Equal(t, "Hello World", echo)
}
