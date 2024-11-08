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

// ***

type ErrKind int

const (
	ErrUnexpected ErrKind = iota + 1
	ErrUnauthorized
	ErrPathNotFound
	ErrValidation
	ErrDataNotFound
	ErrProcessFailed
	ErrRateLimited
)

var RespErrorPresets map[ErrKind]ServerError = map[ErrKind]ServerError{
	ErrUnexpected:    {Code: "E1", Message: "unexpected error"},
	ErrUnauthorized:  {Code: "E2", Message: "not authorized for resource"},
	ErrPathNotFound:  {Code: "E3", Message: "path not found"},
	ErrValidation:    {Code: "E4", Message: "validation error"},
	ErrDataNotFound:  {Code: "E5", Message: "data not found"},
	ErrProcessFailed: {Code: "E6", Message: "process failed"},
	ErrRateLimited:   {Code: "E7", Message: "rate limited"},
}
