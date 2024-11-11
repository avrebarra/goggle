package utils

import (
	"strings"

	"github.com/pkg/errors"
)

type ErrorStackTrace struct {
	FuncName string
	Source   string
}

func ExtractStackTrace(in error) (out []ErrorStackTrace, err error) {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	st, ok := in.(stackTracer)
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
