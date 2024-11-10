package clientgithubrepo

import (
	"context"
	"fmt"
)

var (
	ErrConnectionProblem = fmt.Errorf("connection problem")
)

type Client interface {
	GetPopularRepoNames(ctx context.Context) (names []string, err error)
}
