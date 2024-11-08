package ctxboard

import (
	"context"
	"sync"
)

var refkey = struct{}{}

type ContextBoard struct {
	data *sync.Map
}

func (c ContextBoard) Default() (out ContextBoard) {
	return ContextBoard{data: &sync.Map{}}
}

func CreateWith(in context.Context) (out context.Context) {
	ctx := ContextBoard{}.Default()
	out = context.WithValue(in, refkey, &ctx)
	return
}

func ExtractFrom(ctx context.Context) (ctxdata *ContextBoard, ok bool) {
	if ctxdata, ok = ctx.Value(refkey).(*ContextBoard); !ok {
		ctx := ContextBoard{}.Default()
		ctxdata = &ctx
	}
	return
}
