package rpcserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/avrebarra/goggle/internal/utils"
	"github.com/pkg/errors"

	"github.com/gorilla/rpc"
)

var _ rpc.Codec = (*Codec)(nil)

type Codec struct{}

type ServerRequest struct {
	ID     *json.RawMessage `json:"id"`
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
}

type ServerResponse struct {
	ID     *json.RawMessage `json:"id"`
	Result interface{}      `json:"result,omitempty"`
	Error  interface{}      `json:"error,omitempty"`
}

func (c *Codec) NewRequest(r *http.Request) rpc.CodecRequest {
	req := new(ServerRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	return &CodecRequest{request: req, err: err}
}

type CodecRequest struct {
	request *ServerRequest
	err     error
}

func (c *CodecRequest) Method() (string, error) {
	if c.err == nil {
		return c.request.Method, nil
	}
	return "", c.err
}

func (c *CodecRequest) ReadRequest(args interface{}) error {
	if c.err == nil {
		if c.request.Params != nil {
			params := [1]interface{}{args}
			c.err = json.Unmarshal(*c.request.Params, &params)
		} else {
			c.err = errors.New("rpc: method request ill-formed: missing params field")
		}
	}
	return c.err
}

func (c *CodecRequest) WriteResponse(w http.ResponseWriter, reply interface{}, methodErr error) error {
	if c.err != nil {
		return c.err
	}
	res := &ServerResponse{
		ID:     c.request.ID,
		Result: reply,
		Error:  nil,
	}
	if methodErr != nil {
		_, b, _, _ := runtime.Caller(0)
		basepath := filepath.Dir(b)
		stacktrace, errExtract := utils.ExtractStackTrace(methodErr)
		if errExtract == nil {
			for _, e := range stacktrace {
				fmt.Println(e.FuncName)
				fmt.Println(e.Source)
				fmt.Println()
				if strings.HasPrefix(e.Source, basepath) {
					break
				}
			}
		}

		res.Error = RespErrorPresets[ErrUnexpected].WithMessage(methodErr.Error())
		res.Result = nil
	}
	var serverError ServerError
	if ok := errors.As(methodErr, &serverError); ok {
		res.Error = serverError
		res.Result = nil
	}
	if c.request.ID == nil {
		res.ID = nil
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		encoder := json.NewEncoder(w)
		c.err = encoder.Encode(res)
	}
	return c.err
}
