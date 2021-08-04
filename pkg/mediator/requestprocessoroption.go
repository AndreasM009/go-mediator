package mediator

type RequestProcessorOption func(*Mediator)

func WithPreRequestProcessor(p PreRequestProcessor) RequestProcessorOption {
	return func(m *Mediator) {
		m.preRequestProcessors = append(m.preRequestProcessors, p)
	}
}

func WithPostRequestProcessor(p PostRequestProcessor) RequestProcessorOption {
	return func(m *Mediator) {
		m.postRequestProcessors = append(m.postRequestProcessors, p)
	}
}
