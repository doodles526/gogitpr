package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Repo interface {
	Get(args *RepoArgs) ([]RepoData, error)
}

type repo struct {
	g *ghAPI
}

type RepoArgs struct {
	User string
	Org  string
}

func (r *RepoArgs) validate() error {
	if len(r.User) != 0 && len(r.Org) == 0 {
		return UserOrgErr
	} else if len(r.User) == 0 && len(r.Org) == 0 {
		return UserOrgErr
	}

	return nil
}

func (r *repo) Get(args *RepoArgs) ([]RepoData, error) {
	if err := args.validate(); err != nil {
		return nil, err
	}

	repos := make([]RepoData, 0)
	reqArgs := r.formRequestArgs(args.User, args.Org)

	if err := r.g.doFullPagination(reqArgs, extractRepos(&repos)); err != nil {
		return nil, err
	}

	return repos, nil
}

func (r *repo) formRequestArgs(user, org string) *requestArgs {
	var endpoint string
	if len(user) != 0 {
		endpoint = fmt.Sprintf("/users/%s/repos", user)
	} else {
		endpoint = fmt.Sprintf("/orgs/%s/repos", org)
	}

	return &requestArgs{
		endpoint: endpoint,
		method:   "GET",
	}
}

func extractRepos(repos *[]RepoData) processFunc {
	return func(resp *http.Response) error {
		defer resp.Body.Close()

		repoTmp := make([]RepoData, 0)

		decoder := json.NewDecoder(resp.Body)

		if err := decoder.Decode(&repoTmp); err != nil {
			return err
		}
		*repos = append(*repos, repoTmp...)

		return nil
	}
}

type RepoData struct {
	ID               int         `json:"id"`
	Owner            UserData    `json:"owner"`
	Name             string      `json:"name"`
	FullName         string      `json:"full_name"`
	Description      string      `json:"description"`
	Private          bool        `json:"private"`
	Fork             bool        `json:"fork"`
	URL              string      `json:"url"`
	HTMLURL          string      `json:"html_url"`
	ArchiveURL       string      `json:"archive_url"`
	AssigneesURL     string      `json:"assignees_url"`
	BlobsURL         string      `json:"blobs_url"`
	BranchesURL      string      `json:"branches_url"`
	CloneURL         string      `json:"clone_url"`
	CollaboratorsURL string      `json:"collaborators_url"`
	CommentsURL      string      `json:"comments_url"`
	CommitsURL       string      `json:"commits_url"`
	CompareURL       string      `json:"compare_url"`
	ContentsURL      string      `json:"contents_url"`
	ContributorsURL  string      `json:"contributors_url"`
	DeploymentsURL   string      `json:"deployments_url"`
	DownloadsURL     string      `json:"downloads_url"`
	EventsURL        string      `json:"events_url"`
	ForksURL         string      `json:"forks_url"`
	GitCommitsURL    string      `json:"git_commits_url"`
	GitRefsURL       string      `json:"git_refs_url"`
	GitTagsURL       string      `json:"git_tags_url"`
	GitURL           string      `json:"git_url"`
	HooksURL         string      `json:"hooks_url"`
	IssueCommentURL  string      `json:"issue_comment_url"`
	IssueEventsURL   string      `json:"issue_events_url"`
	IssuesURL        string      `json:"issues_url"`
	KeysURL          string      `json:"keys_url"`
	LabelsURL        string      `json:"labels_url"`
	LanguagesURL     string      `json:"languages_url"`
	MergesURL        string      `json:"merges_url"`
	MilestonesURL    string      `json:"milestones_url"`
	MirrorURL        string      `json:"mirror_url"`
	NotificationsURL string      `json:"notifications_url"`
	PullsURL         string      `json:"pulls_url"`
	ReleasesURL      string      `json:"releases_url"`
	SSHURL           string      `json:"ssh_url"`
	StargazersURL    string      `json:"stargazers_url"`
	StatusesURL      string      `json:"statuses_url"`
	SubscribersURL   string      `json:"subscribers_url"`
	SubscriptionURL  string      `json:"subscription_url"`
	SvnURL           string      `json:"svn_url"`
	TagsURL          string      `json:"tags_url"`
	TeamsURL         string      `json:"teams_url"`
	TreesURL         string      `json:"trees_url"`
	Homepage         string      `json:"homepage"`
	Language         interface{} `json:"language"`
	ForksCount       int         `json:"forks_count"`
	StargazersCount  int         `json:"stargazers_count"`
	WatchersCount    int         `json:"watchers_count"`
	Size             int         `json:"size"`
	DefaultBranch    string      `json:"default_branch"`
	OpenIssuesCount  int         `json:"open_issues_count"`
	Topics           []string    `json:"topics"`
	HasIssues        bool        `json:"has_issues"`
	HasWiki          bool        `json:"has_wiki"`
	HasPages         bool        `json:"has_pages"`
	HasDownloads     bool        `json:"has_downloads"`
	Archived         bool        `json:"archived"`
	PushedAt         time.Time   `json:"pushed_at"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	Permissions      struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	} `json:"permissions"`
	AllowRebaseMerge bool `json:"allow_rebase_merge"`
	AllowSquashMerge bool `json:"allow_squash_merge"`
	AllowMergeCommit bool `json:"allow_merge_commit"`
	SubscribersCount int  `json:"subscribers_count"`
	NetworkCount     int  `json:"network_count"`
}
