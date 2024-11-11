package clientgithubrepo

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/avrebarra/goggle/utils/validator"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type ConfigHTTP struct {
	RESTClient *resty.Client `validate:"required"`
	BaseURL    string        `validate:"required,endswith=/"`
}

type HTTP struct {
	config ConfigHTTP
}

func NewHTTP(cfg ConfigHTTP) (Client, error) {
	if err := validator.Validate(cfg); err != nil {
		err = errors.Wrapf(err, "bad config")
		return nil, err
	}
	e := &HTTP{config: cfg}
	return e, nil
}

func (e *HTTP) GetTopRepoDetails(ctx context.Context) (out []RepoDetail, err error) {
	type ResponseBody struct {
		TotalCount        int  `json:"total_count"`
		IncompleteResults bool `json:"incomplete_results"`
		Items             []struct {
			ID       int    `json:"id"`
			NodeID   string `json:"node_id"`
			Name     string `json:"name"`
			FullName string `json:"full_name"`
			Private  bool   `json:"private"`
			Owner    struct {
				Login             string `json:"login"`
				ID                int    `json:"id"`
				NodeID            string `json:"node_id"`
				AvatarURL         string `json:"avatar_url"`
				GravatarID        string `json:"gravatar_id"`
				URL               string `json:"url"`
				HTMLURL           string `json:"html_url"`
				FollowersURL      string `json:"followers_url"`
				FollowingURL      string `json:"following_url"`
				GistsURL          string `json:"gists_url"`
				StarredURL        string `json:"starred_url"`
				SubscriptionsURL  string `json:"subscriptions_url"`
				OrganizationsURL  string `json:"organizations_url"`
				ReposURL          string `json:"repos_url"`
				EventsURL         string `json:"events_url"`
				ReceivedEventsURL string `json:"received_events_url"`
				Type              string `json:"type"`
				SiteAdmin         bool   `json:"site_admin"`
			} `json:"owner"`
			HTMLURL          string      `json:"html_url"`
			Description      string      `json:"description"`
			Fork             bool        `json:"fork"`
			URL              string      `json:"url"`
			ForksURL         string      `json:"forks_url"`
			KeysURL          string      `json:"keys_url"`
			CollaboratorsURL string      `json:"collaborators_url"`
			TeamsURL         string      `json:"teams_url"`
			HooksURL         string      `json:"hooks_url"`
			IssueEventsURL   string      `json:"issue_events_url"`
			EventsURL        string      `json:"events_url"`
			AssigneesURL     string      `json:"assignees_url"`
			BranchesURL      string      `json:"branches_url"`
			TagsURL          string      `json:"tags_url"`
			BlobsURL         string      `json:"blobs_url"`
			GitTagsURL       string      `json:"git_tags_url"`
			GitRefsURL       string      `json:"git_refs_url"`
			TreesURL         string      `json:"trees_url"`
			StatusesURL      string      `json:"statuses_url"`
			LanguagesURL     string      `json:"languages_url"`
			StargazersURL    string      `json:"stargazers_url"`
			ContributorsURL  string      `json:"contributors_url"`
			SubscribersURL   string      `json:"subscribers_url"`
			SubscriptionURL  string      `json:"subscription_url"`
			CommitsURL       string      `json:"commits_url"`
			GitCommitsURL    string      `json:"git_commits_url"`
			CommentsURL      string      `json:"comments_url"`
			IssueCommentURL  string      `json:"issue_comment_url"`
			ContentsURL      string      `json:"contents_url"`
			CompareURL       string      `json:"compare_url"`
			MergesURL        string      `json:"merges_url"`
			ArchiveURL       string      `json:"archive_url"`
			DownloadsURL     string      `json:"downloads_url"`
			IssuesURL        string      `json:"issues_url"`
			PullsURL         string      `json:"pulls_url"`
			MilestonesURL    string      `json:"milestones_url"`
			NotificationsURL string      `json:"notifications_url"`
			LabelsURL        string      `json:"labels_url"`
			ReleasesURL      string      `json:"releases_url"`
			DeploymentsURL   string      `json:"deployments_url"`
			CreatedAt        time.Time   `json:"created_at"`
			UpdatedAt        time.Time   `json:"updated_at"`
			PushedAt         time.Time   `json:"pushed_at"`
			GitURL           string      `json:"git_url"`
			SSHURL           string      `json:"ssh_url"`
			CloneURL         string      `json:"clone_url"`
			SvnURL           string      `json:"svn_url"`
			Homepage         string      `json:"homepage"`
			Size             int         `json:"size"`
			StargazersCount  int         `json:"stargazers_count"`
			WatchersCount    int         `json:"watchers_count"`
			Language         string      `json:"language"`
			HasIssues        bool        `json:"has_issues"`
			HasProjects      bool        `json:"has_projects"`
			HasDownloads     bool        `json:"has_downloads"`
			HasWiki          bool        `json:"has_wiki"`
			HasPages         bool        `json:"has_pages"`
			ForksCount       int         `json:"forks_count"`
			MirrorURL        interface{} `json:"mirror_url"`
			Archived         bool        `json:"archived"`
			Disabled         bool        `json:"disabled"`
			OpenIssuesCount  int         `json:"open_issues_count"`
			License          struct {
				Key    string `json:"key"`
				Name   string `json:"name"`
				SpdxID string `json:"spdx_id"`
				URL    string `json:"url"`
				NodeID string `json:"node_id"`
			} `json:"license"`
			Forks         int     `json:"forks"`
			OpenIssues    int     `json:"open_issues"`
			Watchers      int     `json:"watchers"`
			DefaultBranch string  `json:"default_branch"`
			Score         float64 `json:"score"`
		} `json:"items"`
	}

	// ***

	// prepare request
	data := url.Values{}
	data.Set("sort", "stars")
	data.Set("order", "desc")
	data.Set("q", "language:go")
	data.Set("q", fmt.Sprintf("created:>%s", time.Now().Add(-1*7*24*time.Hour).Format("2006-01-02")))

	respdata := ResponseBody{}
	resp, err := e.config.RESTClient.R().
		SetContext(ctx).
		SetResult(&respdata).
		Get(e.config.BaseURL + "search/repositories" + "?" + data.Encode())
	if err != nil {
		err = fmt.Errorf("http request failure: %w", err)
		return
	}
	if resp.StatusCode() != 200 {
		err = fmt.Errorf("http request failure: bad http status code: %s", resp.Status())
		return
	}

	out = []RepoDetail{}
	for _, e := range respdata.Items {
		out = append(out, RepoDetail{
			Name:            e.Name,
			StargazersCount: e.StargazersCount,
			Author:          e.Owner.Login,
			URI:             e.HTMLURL,
		})
	}

	return
}
