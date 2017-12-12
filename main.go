package main

import (
	"github.com/doodles526/gogitpr/api"
	"github.com/doodles526/gogitpr/config"
	"github.com/doodles526/gogitpr/db"

	"fmt"
	"os"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Printf("Error creating config: %+v", err)
		os.Exit(1)
	}

	apiArgs := &api.GithubAPIArgs{
		BaseURL:         cfg.BaseURL,
		Token:           cfg.GithubToken,
		ApplicationName: cfg.ApplicationName,
		Logger:          cfg.Logger,
		// use default version
	}

	gh, err := api.NewGithubAPI(apiArgs)
	if err != nil {
		fmt.Printf("Error creating GithubAPI: %+v", err)
		os.Exit(1)
	}

	prArgs := &api.PullRequestArgs{
		User: cfg.GithubUser,
		Org:  cfg.GithubOrg,
	}

	prs, err := gh.PullRequest().Get(prArgs)
	if err != nil {
		fmt.Printf("Error getting Pull Requests: %+v", err)
		os.Exit(1)
	}

	dbArgs := &db.Args{
		Logger: cfg.Logger,
	}

	prDB, err := db.NewDB(dbArgs)
	if err != nil {
		fmt.Printf("Error Creating DB: %+v", err)
		os.Exit(1)
	}

	if err := prDB.StorePullRequestBatch(prs); err != nil {
		fmt.Printf("Error storing PRs: %+v", err)
		os.Exit(1)
	}

	allPRs, err := prDB.GetAllPullRequests()
	if err != nil {
		fmt.Printf("Error Fetching PRs: %+v", err)
		os.Exit(1)
	}

	if cfg.PrintResult {
		fmt.Println(allPRs)
	}
}
