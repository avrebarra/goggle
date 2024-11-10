package clientgithubrepo

import "context"

type Dummy struct {
}

func (Dummy) GetPopularRepoNames(ctx context.Context) (names []string, err error) {
	names = []string{"sample"}
	return
}
