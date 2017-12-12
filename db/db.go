package db

import (
	"github.com/doodles526/gogitpr/api"
	"github.com/sirupsen/logrus"
)

// PRFilterFunc should return true if we should return this particular PR
// Only return an error if unrecoverable
type PRFilterFunc func(pr api.PullRequestData) (bool, error)

// DB is an abstraction on top of whatever backing store exists
// for purposes of interview brevity this will just be an in-mem
// storage. But easily extendable to a persisted DB like RethinkDB
type DB interface {
	StorePullRequest(pr api.PullRequestData) error
	StorePullRequestBatch(prs []api.PullRequestData) error
	GetAllPullRequests() ([]api.PullRequestData, error)
	GetFilterPullRequests(f PRFilterFunc) ([]api.PullRequestData, error)
	GetPullRequestByID(id int) (api.PullRequestData, bool, error)
}

// DBArgs is currently empty, as to be forward compatible
// if we ever choose to back by persistent storage
type DBArgs struct {
	Logger *logrus.Logger
}

func NewDB(args *DBArgs) (DB, error) {
	return &inMem{
		pullRequests: make([]api.PullRequestData, 0),
		idIndex:      make(map[int]*api.PullRequestData),
		logger:       args.Logger.WithFields(logrus.Fields{"prefix": "DB"}),
	}, nil
}

// inMem is not concurrency safe. The client must perform locking if they wish
// to have safe concurrent read/write access
type inMem struct {
	pullRequests []api.PullRequestData

	idIndex map[int]*api.PullRequestData
	logger  *logrus.Entry
}

func (i *inMem) StorePullRequest(pr api.PullRequestData) error {
	i.pullRequests = append(i.pullRequests, pr)

	i.idIndex[pr.ID] = &i.pullRequests[len(i.pullRequests)-1]

	return nil
}

func (i *inMem) StorePullRequestBatch(prs []api.PullRequestData) error {
	for _, pr := range prs {
		if err := i.StorePullRequest(pr); err != nil {
			return err
		}
	}

	return nil
}

func (i *inMem) GetAllPullRequests() ([]api.PullRequestData, error) {
	prTemp := make([]api.PullRequestData, len(i.pullRequests))
	// copy so we don't pass the backing store's copy of the slice
	copy(prTemp, i.pullRequests)

	return prTemp, nil
}

func (i *inMem) GetFilterPullRequests(f PRFilterFunc) ([]api.PullRequestData, error) {
	prTemp := make([]api.PullRequestData, 0)
	for _, pr := range i.pullRequests {
		ok, err := f(pr)
		if err != nil {
			return nil, err
		}
		if ok {
			prTemp = append(prTemp, pr)
		}
	}

	return prTemp, nil
}

func (i *inMem) GetPullRequestByID(id int) (api.PullRequestData, bool, error) {
	pr, ok := i.idIndex[id]
	if !ok {
		return api.PullRequestData{}, false, nil
	}
	newPR := new(api.PullRequestData)
	*newPR = *pr

	return *newPR, true, nil
}
