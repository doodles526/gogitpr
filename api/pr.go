package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PullRequestArgs struct {
	User  string
	Org   string
	Repos []string
}

func (a *PullRequestArgs) validate() error {
	if len(a.User) != 0 && len(a.Org) != 0 {
		return UserOrgErr
	}

	if len(a.User) == 0 && len(a.Org) == 0 {
		return UserOrgErr
	}

	return nil
}

type PullRequest interface {
	Get(args *PullRequestArgs) ([]PullRequestData, error)
}

type pullRequest struct {
	g *ghAPI
}

func (p *pullRequest) reqString() string {
	return fmt.Sprintf("")
}

func (p *pullRequest) Get(args *PullRequestArgs) ([]PullRequestData, error) {
	if err := args.validate(); err != nil {
		return nil, err
	}

	if len(args.Repos) == 0 {
		if err := p.populateRepos(args); err != nil {
			return nil, err
		}
	}

	pullRequests := make([]PullRequestData, 0)
	for _, repo := range args.Repos {
		reqArgs := p.formRequestArgs(args.User, args.Org, repo)

		if err := p.g.doFullPagination(reqArgs, extractPRs(&pullRequests)); err != nil {
			return nil, err
		}
	}

	return pullRequests, nil
}

func extractPRs(prData *[]PullRequestData) processFunc {
	return func(resp *http.Response) error {
		defer resp.Body.Close()

		prTmp := make([]PullRequestData, 0)

		decoder := json.NewDecoder(resp.Body)

		if err := decoder.Decode(&prTmp); err != nil {
			return err
		}
		*prData = append(*prData, prTmp...)

		return nil
	}
}

func (p *pullRequest) formRequestArgs(user, org, repo string) *requestArgs {
	var owner string
	if len(user) != 0 {
		owner = user
	} else {
		owner = org
	}
	endpoint := fmt.Sprintf("/repos/%s/%s/pulls", owner, repo)

	return &requestArgs{
		endpoint: endpoint,
		method:   "GET",
	}
}

func (p *pullRequest) populateRepos(args *PullRequestArgs) error {
	var repoArgs *RepoArgs
	if len(args.User) != 0 {
		repoArgs = &RepoArgs{
			User: args.User,
		}
	} else if len(args.Org) != 0 {
		repoArgs = &RepoArgs{
			Org: args.Org,
		}
	} else {
		repoArgs = &RepoArgs{}
	}

	repos, err := p.g.Repos().Get(repoArgs)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		args.Repos = append(args.Repos, repo.Name)
	}

	return nil
}

// PullRequestData auto-generated from https://mholt.github.io/json-to-go/
type PullRequestData struct {
	ID                int           `json:"id"`
	URL               string        `json:"url"`
	HTMLURL           string        `json:"html_url"`
	DiffURL           string        `json:"diff_url"`
	PatchURL          string        `json:"patch_url"`
	IssueURL          string        `json:"issue_url"`
	CommitsURL        string        `json:"commits_url"`
	ReviewCommentsURL string        `json:"review_comments_url"`
	ReviewCommentURL  string        `json:"review_comment_url"`
	CommentsURL       string        `json:"comments_url"`
	StatusesURL       string        `json:"statuses_url"`
	Number            int           `json:"number"`
	State             string        `json:"state"`
	Title             string        `json:"title"`
	Body              string        `json:"body"`
	Assignee          UserData      `json:"assignee"`
	Milestone         MilestoneData `json:"milestone"`
	Locked            bool          `json:"locked"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
	ClosedAt          time.Time     `json:"closed_at"`
	MergedAt          time.Time     `json:"merged_at"`
	Head              CommitData    `json:"head"`
	Base              CommitData    `json:"base"`
	Links             struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
		Issue struct {
			Href string `json:"href"`
		} `json:"issue"`
		Comments struct {
			Href string `json:"href"`
		} `json:"comments"`
		ReviewComments struct {
			Href string `json:"href"`
		} `json:"review_comments"`
		ReviewComment struct {
			Href string `json:"href"`
		} `json:"review_comment"`
		Commits struct {
			Href string `json:"href"`
		} `json:"commits"`
		Statuses struct {
			Href string `json:"href"`
		} `json:"statuses"`
	} `json:"_links"`
	User UserData `json:"user"`
}
