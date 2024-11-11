package utils

import (
	"strings"

	"github.com/pkg/errors"
)

type ErrorStackTrace struct {
	FuncName string
	Source   string
}

func ExtractStackTrace(cur error) (out []ErrorStackTrace, err error) {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	istracer := func(cur error) bool {
		if cur == nil {
			return false
		}
		_, ok := cur.(stackTracer)
		return ok
	}
	nexttracer := func(cur error) error {
		for cur != nil {
			child := errors.Unwrap(cur)
			if istracer(child) {
				return child
			}
			cur = child
		}
		return nil
	}

	for {
		child := nexttracer(cur)
		if child == nil {
			break
		}
		cur = child
	}

	st, ok := cur.(stackTracer)
	if !ok {
		err = errors.New("err have no stacktrace")
		return
	}

	out = []ErrorStackTrace{}
	for _, f := range st.StackTrace() {
		componentsraw, _ := f.MarshalText()
		components := strings.SplitN(string(componentsraw), " ", 2)
		fx, pat := components[0], components[1]
		out = append(out, ErrorStackTrace{
			FuncName: fx,
			Source:   pat,
		})
	}

	return
}
