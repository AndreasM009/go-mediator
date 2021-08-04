package mediator

type Response interface {
	Result() interface{}
	HasError() bool
	Error() error
}

type response struct {
	result interface{}
	err    error
}

type emptyResponse struct {
	err error
}

func CreateResponse(result interface{}, err error) Response {
	return &response{
		result: result,
		err:    err,
	}
}

func CreateEmtpyResponse(err error) Response {
	return &emptyResponse{
		err: err,
	}
}

func (r *response) Result() interface{} {
	return r.result
}

func (r *response) HasError() bool {
	return r.err != nil
}

func (r *response) Error() error {
	return r.err
}

func (er *emptyResponse) Result() interface{} {
	return struct{}{}
}

func (er *emptyResponse) HasError() bool {
	return nil != er.err
}

func (er *emptyResponse) Error() error {
	return er.err
}
