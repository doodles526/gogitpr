package db

import (
	"testing"

	"github.com/doodles526/gogitpr/api"
	"github.com/stretchr/testify/assert"
)

func TestStorePullRequest(t *testing.T) {
	db := &inMem{
		pullRequests: make([]api.PullRequestData, 0),
		idIndex:      make(map[int]*api.PullRequestData),
	}

	pr := api.PullRequestData{
		ID: 1234,
	}

	err := db.StorePullRequest(pr)
	assert.NoError(t, err, "Should be no error storing PR")

	lenDB := len(db.pullRequests)

	assert.Equal(t, lenDB, 1, "Length should be 1")

	pr2 := db.pullRequests[0]

	assert.Equal(t, pr2.ID, pr.ID, "Should have matching values")
}

func TestGetAllPullRequests(t *testing.T) {
	pr := api.PullRequestData{
		ID: 1234,
	}
	db := &inMem{
		pullRequests: []api.PullRequestData{pr},
		idIndex:      make(map[int]*api.PullRequestData),
	}

	prs, err := db.GetAllPullRequests()

	assert.NoError(t, err, "Should be no error getting Pull requests")

	lenPRs := len(prs)
	assert.Equal(t, lenPRs, 1, "Should get back 1 PR")

	assert.Equal(t, prs[0].ID, pr.ID, "Should get the proper name back")

	pr.ID = 4321

	assert.NotEqual(t, prs[0].ID, pr.ID, "Changing our local version should not change DBs")
}

func TestGetFilterPullRequests(t *testing.T) {
	pr := api.PullRequestData{
		ID: 1234,
	}
	pr2 := api.PullRequestData{
		ID: 4321,
	}
	db := &inMem{
		pullRequests: []api.PullRequestData{pr, pr2},
		idIndex:      make(map[int]*api.PullRequestData),
	}

	filterFunc := func(pr api.PullRequestData) (bool, error) {
		return pr.ID == 1234, nil
	}

	prs, err := db.GetFilterPullRequests(filterFunc)

	assert.NoError(t, err, "Should have no error getting prs")

	lenPR := len(prs)

	assert.Equal(t, 1, lenPR, "Should have gotten 1 PR")
	assert.Equal(t, 1234, prs[0].ID, "Should have gotten the PR from our filter function")
}

func TestGetPullRequestByID(t *testing.T) {
	pr := api.PullRequestData{
		ID: 1234,
	}
	pr2 := api.PullRequestData{
		ID: 4321,
	}
	db := &inMem{
		pullRequests: []api.PullRequestData{pr, pr2},
		idIndex: map[int]*api.PullRequestData{
			1234: &pr,
			4321: &pr2,
		},
	}

	prBack, ok, err := db.GetPullRequestByID(1234)
	assert.NoError(t, err, "Should be no error fetching PR")
	assert.True(t, ok, "Should get back a pr")
	assert.Equal(t, 1234, prBack.ID, "Should have the correct PR back")

	_, ok, err = db.GetPullRequestByID(8484)
	assert.NoError(t, err, "Shuold be no error looking for PR")
	assert.False(t, ok, "Should not be a PR back")
}
