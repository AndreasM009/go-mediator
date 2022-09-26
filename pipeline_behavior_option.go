package mediator

type PipelineBehaviorOption struct {
	processor PipelineBehavior
}

func WithBehavior(p PipelineBehavior) PipelineBehaviorOption {
	return PipelineBehaviorOption{
		processor: p,
	}
}
