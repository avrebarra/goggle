package rpcserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/avrebarra/goggle/utils/ctxboard"
	"github.com/gorilla/rpc"
	"github.com/pkg/errors"
)

var _ rpc.Codec = (*Codec)(nil)

type Codec struct{}

type ServerRequest struct {
	ID     string           `json:"id"`
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
}

type ServerResponse struct {
	ID     string      `json:"id"`
	Result interface{} `json:"result,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}

func (c *Codec) NewRequest(r *http.Request) rpc.CodecRequest {
	req := new(ServerRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		panic(err)
	}

	ctx := r.Context()
	reqctx := ctxboard.GetData(ctx, KeyRequestContext).(*RequestContext)
	reqctx.OpsID = req.ID
	reqctx.OpsName = req.Method
	reqctx.IngoingData = req.Params

	return &CodecRequest{ctx: r.Context(), request: req}
}

type CodecRequest struct {
	ctx     context.Context
	request *ServerRequest
}

func (c *CodecRequest) Method() (string, error) {
	return c.request.Method, nil
}

func (c *CodecRequest) ReadRequest(args interface{}) (err error) {
	if c.request.Params != nil {
		params := [1]interface{}{args}
		err = json.Unmarshal(*c.request.Params, &params)
	} else {
		err = errors.New("rpc: method request ill-formed: missing params field")
	}

	return
}

func (c *CodecRequest) WriteResponse(w http.ResponseWriter, reply interface{}, err error) error {
	if err != nil {
		panic(err)
	}

	ctx := c.ctx
	reqctx := ctxboard.GetData(ctx, KeyRequestContext).(*RequestContext)
	reqctx.OutgoingData = reply

	res := &ServerResponse{
		ID:     c.request.ID,
		Result: reply,
		Error:  nil,
	}

	if c.request.ID == "" {
		res.ID = ""
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		encoder := json.NewEncoder(w)
		err = encoder.Encode(res)
	}
	if err != nil {
		panic(err)
	}

	return nil
}
