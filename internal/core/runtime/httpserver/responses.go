package httpserver

var (
	ErrUnexpected = ServerError{Code: "unexpected", Cause: "unexpected error"}
	ErrNotFound   = ServerError{Code: "not_found", Cause: "resource not found"}
	ErrValidation = ServerError{Code: "invalid", Cause: "validation error"}
	ErrDuplicate  = ServerError{Code: "duplicate", Cause: "resource already exists"}
)

const (
	RespSuccess RespKind = iota
	RespUnauthorized
	RespNotFound
	RespProcessFailed
	RespBadRequest
	RespUnexpected
)

var respmapper map[RespKind]preset = map[RespKind]preset{
	RespSuccess:       {Status: 200, Resp: Resp{Message: "Request processed successfully!"}},
	RespBadRequest:    {Status: 400, Resp: Resp{Message: "Your request was not valid. Please try again with valid data."}},
	RespUnauthorized:  {Status: 401, Resp: Resp{Message: "You are not authorized to access specified resource."}},
	RespNotFound:      {Status: 404, Resp: Resp{Message: "Cannot find requested resource on server."}},
	RespProcessFailed: {Status: 422, Resp: Resp{Message: "Internal process failed. Please check your request input and try again."}},
	RespUnexpected:    {Status: 500, Resp: Resp{Message: "Encountered error while processing request."}},
}
