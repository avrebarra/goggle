package rpcserver

import (
	"net/http"
	"time"

	"github.com/jinzhu/copier"
)

func (s *Handler) Ping(r *http.Request, in *ReqPing, out *RespPing) (err error) {
	resp := RespPing{
		Version:   s.Version,
		StartedAt: s.StartedAt,
		Uptime:    time.Since(s.StartedAt).Round(time.Second).String(),
	}
	return copier.Copy(out, resp)
}
