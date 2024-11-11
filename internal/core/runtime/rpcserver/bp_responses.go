package rpcserver

type ServerError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e ServerError) Error() string {
	return e.Message
}

func (e ServerError) WithMessage(msg string) ServerError {
	e.Message += ": " + msg
	return e
}
