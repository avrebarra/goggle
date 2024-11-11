package clientgithubrepo

import (
	"context"
	"fmt"
)

var (
	ErrConnectionProblem = fmt.Errorf("connection problem")
)

type Client interface {
	GetTopRepoDetails(ctx context.Context) (out []RepoDetail, err error)
}

type RepoDetail struct {
	Name            string
	Author          string
	URI             string
	StargazersCount int
}
