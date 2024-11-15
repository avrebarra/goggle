package clientgithubrepo_test

import (
	"context"
	"testing"

	"github.com/avrebarra/goggle/internal/module/clientgithubrepo"
	"github.com/avrebarra/goggle/internal/utils"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	HTTP_DefaultConfig = clientgithubrepo.ConfigHTTP{
		RESTClient: resty.New(),
		BaseURL:    "https://api.github.com/",
	}
)

func TestNew(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		c, err := clientgithubrepo.NewHTTP(clientgithubrepo.ConfigHTTP{
			RESTClient: resty.New(),
			BaseURL:    "https://api.github.com/",
		})

		assert.NoError(t, err)
		assert.NotNil(t, c)
	})

	t.Run("err bad dep", func(t *testing.T) {
		_, err := clientgithubrepo.NewHTTP(clientgithubrepo.ConfigHTTP{})

		assert.Error(t, err)
	})
}

func TestHTTP_GetPopularRepoNames(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		cfg := HTTP_DefaultConfig

		httpmock.ActivateNonDefault(cfg.RESTClient.GetClient())
		defer httpmock.DeactivateAndReset()

		fixture := `{"total_count":675762,"incomplete_results":false,"items":[{"id":396474980,"node_id":"MDEwOlJlcG9zaXRvcnkzOTY0NzQ5ODA=","name":"AppleNeuralHash2ONNX","full_name":"AsuharietYgvar/AppleNeuralHash2ONNX","private":false,"owner":{"login":"AsuharietYgvar","id":88948101,"node_id":"MDQ6VXNlcjg4OTQ4MTAx","avatar_url":"https://avatars.githubusercontent.com/u/88948101?v=4","gravatar_id":"","url":"https://api.github.com/users/AsuharietYgvar","html_url":"https://github.com/AsuharietYgvar","followers_url":"https://api.github.com/users/AsuharietYgvar/followers","following_url":"https://api.github.com/users/AsuharietYgvar/following{/other_user}","gists_url":"https://api.github.com/users/AsuharietYgvar/gists{/gist_id}","starred_url":"https://api.github.com/users/AsuharietYgvar/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/AsuharietYgvar/subscriptions","organizations_url":"https://api.github.com/users/AsuharietYgvar/orgs","repos_url":"https://api.github.com/users/AsuharietYgvar/repos","events_url":"https://api.github.com/users/AsuharietYgvar/events{/privacy}","received_events_url":"https://api.github.com/users/AsuharietYgvar/received_events","type":"User","site_admin":false},"html_url":"https://github.com/AsuharietYgvar/AppleNeuralHash2ONNX","description":"Convert Apple NeuralHash model for CSAM Detection to ONNX.","fork":false,"url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX","forks_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/forks","keys_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/keys{/key_id}","collaborators_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/collaborators{/collaborator}","teams_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/teams","hooks_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/hooks","issue_events_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/issues/events{/number}","events_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/events","assignees_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/assignees{/user}","branches_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/branches{/branch}","tags_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/tags","blobs_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/git/blobs{/sha}","git_tags_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/git/tags{/sha}","git_refs_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/git/refs{/sha}","trees_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/git/trees{/sha}","statuses_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/statuses/{sha}","languages_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/languages","stargazers_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/stargazers","contributors_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/contributors","subscribers_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/subscribers","subscription_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/subscription","commits_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/commits{/sha}","git_commits_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/git/commits{/sha}","comments_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/comments{/number}","issue_comment_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/issues/comments{/number}","contents_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/contents/{+path}","compare_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/compare/{base}...{head}","merges_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/merges","archive_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/{archive_format}{/ref}","downloads_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/downloads","issues_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/issues{/number}","pulls_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/pulls{/number}","milestones_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/milestones{/number}","notifications_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/notifications{?since,all,participating}","labels_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/labels{/name}","releases_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/releases{/id}","deployments_url":"https://api.github.com/repos/AsuharietYgvar/AppleNeuralHash2ONNX/deployments","created_at":"2021-08-15T19:52:47Z","updated_at":"2021-08-20T03:47:17Z","pushed_at":"2021-08-19T18:33:53Z","git_url":"git://github.com/AsuharietYgvar/AppleNeuralHash2ONNX.git","ssh_url":"git@github.com:AsuharietYgvar/AppleNeuralHash2ONNX.git","clone_url":"https://github.com/AsuharietYgvar/AppleNeuralHash2ONNX.git","svn_url":"https://github.com/AsuharietYgvar/AppleNeuralHash2ONNX","homepage":"","size":15,"stargazers_count":966,"watchers_count":966,"language":"Python","has_issues":true,"has_projects":true,"has_downloads":true,"has_wiki":true,"has_pages":false,"forks_count":76,"mirror_url":null,"archived":false,"disabled":false,"open_issues_count":3,"license":{"key":"apache-2.0","name":"Apache License 2.0","spdx_id":"Apache-2.0","url":"https://api.github.com/licenses/apache-2.0","node_id":"MDc6TGljZW5zZTI="},"forks":76,"open_issues":3,"watchers":966,"default_branch":"master","score":1}]}`
		httpmock.RegisterResponder(
			"GET", cfg.BaseURL+"search/repositories",
			httpmock.NewJsonResponderOrPanic(200, utils.UnmarshalToMap([]byte(fixture))),
		)

		e, err := clientgithubrepo.NewHTTP(HTTP_DefaultConfig)
		require.NoError(t, err)

		out, err := e.GetTopRepoDetails(context.Background())

		assert.NoError(t, err)
		assert.NotEmpty(t, out)
	})

	t.Run("err bad http code", func(t *testing.T) {
		cfg := HTTP_DefaultConfig

		httpmock.ActivateNonDefault(cfg.RESTClient.GetClient())
		defer httpmock.DeactivateAndReset()

		fixture := ""
		httpmock.RegisterResponder(
			"GET", cfg.BaseURL+"search/repositories",
			httpmock.NewJsonResponderOrPanic(504, utils.UnmarshalToMap([]byte(fixture))),
		)

		e, err := clientgithubrepo.NewHTTP(HTTP_DefaultConfig)
		require.NoError(t, err)

		_, err = e.GetTopRepoDetails(context.Background())

		assert.Error(t, err)
	})

	t.Run("err bad request", func(t *testing.T) {
		cfg := HTTP_DefaultConfig

		httpmock.ActivateNonDefault(cfg.RESTClient.GetClient())
		defer httpmock.DeactivateAndReset()

		e, err := clientgithubrepo.NewHTTP(HTTP_DefaultConfig)
		require.NoError(t, err)

		_, err = e.GetTopRepoDetails(context.Background())

		assert.Error(t, err)
	})
}
