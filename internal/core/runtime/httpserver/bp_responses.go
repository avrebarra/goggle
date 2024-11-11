package httpserver

type ServerOutput struct {
	Msg   string `json:"message"`
	Data  any    `json:"data,omitempty"`
	Error any    `json:"error,omitempty"`
}

func (e ServerOutput) Normalize() ServerOutput {
	if e.Data == nil {
		e.Data = map[string]interface{}{}
	}
	if e.Error != nil {
		e.Data = nil
	}
	return e
}

type ServerError struct {
	error
	Code   string `json:"code"`
	Cause  string `json:"cause"`
	Detail any    `json:"detail,omitempty"`
}

func (e ServerError) Error() string {
	return e.Cause
}

func (e ServerError) Unwrap() error {
	return e.error
}

func (e ServerError) WithMessage(msg string) ServerError {
	e.Cause += ": " + msg
	return e
}

func (e ServerError) WithDetail(data any) ServerError {
	e.Detail = data
	return e
}

func (e ServerError) Wrap(err error) ServerError {
	e.error = err
	return e
}

// ***

type preset struct {
	Status int
	Resp   Resp
}

// ***

type RespKind int

type Resp struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func (e Resp) Normalize() Resp {
	if e.Data == nil {
		e.Data = map[string]interface{}{}
	}
	if e.Error != nil {
		e.Data = nil
	}
	return e
}
