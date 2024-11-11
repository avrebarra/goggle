package rpcserver

type ErrKind int

const (
	ErrUnexpected ErrKind = iota + 1
	ErrUnauthorized
	ErrNotFound
	ErrValidation
	ErrProcessFailed
	ErrRateLimited
)

var RespErrorPresets map[ErrKind]ServerError = map[ErrKind]ServerError{
	ErrUnexpected:    {Code: "unexpected", Message: "unexpected error"},
	ErrUnauthorized:  {Code: "unauthorized", Message: "not authorized for resource"},
	ErrNotFound:      {Code: "not_found", Message: "resource not found"},
	ErrValidation:    {Code: "invalid", Message: "validation error"},
	ErrProcessFailed: {Code: "failed", Message: "process failed"},
	ErrRateLimited:   {Code: "rate_limited", Message: "rate limited"},
}
