package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/peterhellberg/link" //RFC5988 complient header parser
	"github.com/sirupsen/logrus"
)

const defaultBase = "https://api.github.com"
const defaultVersion = Version3

// Version defines which version of the github API to use
// Currently we only support Version3
type Version int

const (
	// VersionDefault allows us to set a default value if none is specified
	VersionDefault Version = iota
	// Version3 is the github API version 3
	Version3
	// Version4 is the github API version 4
	Version4
)

// GithubAPI provides access to the github API
// TODO: Implement remaining endpoints
type GithubAPI interface {
	PullRequest() PullRequest
	Repos() Repo
}

type ghAPI struct {
	baseURL   *url.URL
	token     string
	userAgent string
	version   Version

	client *http.Client
	logger *logrus.Entry
}

// GithubAPIArgs specifies how the github API should be queried
type GithubAPIArgs struct {
	BaseURL         string
	Token           string
	ApplicationName string
	Version         Version
	Logger          *logrus.Logger
}

// NewGithubAPI creates a new client for accessing the github api
func NewGithubAPI(args *GithubAPIArgs) (GithubAPI, error) {
	if err := args.validate(); err != nil {
		return nil, err
	}

	bu, err := url.Parse(args.BaseURL)
	if err != nil {
		return nil, err
	}

	base := &ghAPI{
		baseURL:   bu,
		token:     args.Token,
		userAgent: args.ApplicationName,
		version:   args.Version,
		client:    &http.Client{},
		logger:    args.Logger.WithFields(logrus.Fields{"prefix": "GithubAPI"}),
	}

	return base, nil
}

func (a *GithubAPIArgs) validate() error {
	if len(a.BaseURL) == 0 {
		a.BaseURL = defaultBase
	}

	if a.Version == VersionDefault {
		a.Version = defaultVersion
	} else if a.Version != Version3 {
		// TODO: Remove this conditional upon V4 support
		return argUnsupported("Version", a.Version)
	}

	if len(a.ApplicationName) == 0 {
		return argMissingError("ApplicationName")
	}

	return nil
}

func (g *ghAPI) PullRequest() PullRequest {
	return &pullRequest{
		g: g,
	}
}

func (g *ghAPI) Repos() Repo {
	return &repo{
		g: g,
	}
}

type requestArgs struct {
	values   map[string]string
	endpoint string
	method   string
}

func deepCopyURL(u *url.URL) *url.URL {
	eNew := new(url.URL)
	*eNew = *u
	return eNew
}

func deepCopyRequestArgs(a *requestArgs) *requestArgs {
	aNew := new(requestArgs)
	*aNew = *a
	return aNew
}

// doRequest performs data request from args
func (g *ghAPI) doRequest(args *requestArgs) (*http.Response, error) {
	u := deepCopyURL(g.baseURL)

	g.logger.Debugf("performing request - %s %s", args.method, args.endpoint)

	u.Path = fmt.Sprintf("%s%s", u.Path, args.endpoint)

	if args.values != nil {
		for key, val := range args.values {
			u.Query().Set(key, val)
		}
	}

	req, err := http.NewRequest(args.method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if len(g.token) != 0 {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", g.token))
	}

	return g.client.Do(req)
}

type processFunc func(*http.Response) error

func (g *ghAPI) doFullPagination(args *requestArgs, f processFunc) error {
	nArgs := deepCopyRequestArgs(args)

	if nArgs.values == nil {
		nArgs.values = make(map[string]string)
	}

	resp, err := g.doRequest(nArgs)
	if err != nil {
		return err
	}

	totalPages := 1

	// We have to process this now since passing it to f() could alter the header
	// and we need to ensure someone isn't able to alter the Link header before this
	if linkGroup := link.ParseHeader(resp.Header); linkGroup != nil {
		if last, ok := linkGroup["last"]; ok {
			linkURL, err := url.Parse(last.URI)
			if err != nil {
				return err
			}
			if pageOfLast := linkURL.Query().Get("page"); pageOfLast != "" {
				iPageOfLast, err := strconv.Atoi(pageOfLast)
				if err != nil {
					return err
				}
				totalPages = iPageOfLast
			}
		}
	}

	if err := f(resp); err != nil {
		return err
	}

	for curPage := 2; curPage <= totalPages; curPage++ {
		nArgs.values["page"] = strconv.Itoa(curPage)

		resp, err := g.doRequest(nArgs)
		if err != nil {
			return err
		}

		if err := f(resp); err != nil {
			return err
		}
	}

	return nil
}

func argMissingError(field string) error {
	return fmt.Errorf("%s must be set in GithubAPIArgs", field)
}

func argUnsupported(field string, value interface{}) error {
	return fmt.Errorf("Currently the value %v is not supported for %s", value, field)
}
