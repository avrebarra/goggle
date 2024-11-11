package cronworker

import (
	"context"
	"time"

	"github.com/avrebarra/goggle/internal/module/clientgithubrepo"
)

func (h *CronHandler) ExecGitHubPing(ctx context.Context) (out any, err error) {
	type OutputResp struct {
		Repos     []clientgithubrepo.RepoDetail `json:"repositories"`
		CheckedAt time.Time                     `json:"checked_at"`
	}

	// ***

	timestamp := time.Now()
	resp, err := h.GithubClient.GetTopRepoDetails(ctx)
	if err != nil {
		return
	}

	out = OutputResp{
		Repos:     resp,
		CheckedAt: timestamp,
	}
	return
}
