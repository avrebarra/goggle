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
	ErrUnexpected:    {Code: "unexpected", Message: "Unexpected error occurred."},
	ErrUnauthorized:  {Code: "unauthorized", Message: "You are not authorized for this resource."},
	ErrNotFound:      {Code: "not_found", Message: "Resource was not found."},
	ErrValidation:    {Code: "invalid", Message: "Request was invalid."},
	ErrProcessFailed: {Code: "failed", Message: "Process failed."},
	ErrRateLimited:   {Code: "rate_limited", Message: "Too many requests, try again later."},
}
