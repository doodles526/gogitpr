package api

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeepCopyURL(t *testing.T) {
	u := &url.URL{
		Host: "localhost:8080",
	}

	u2 := deepCopyURL(u)

	u2.Host = "test"

	assert.NotEqual(t, u.Host, u2.Host)
}

func TestDeepCopyRequestArgs(t *testing.T) {
	a := &requestArgs{
		method: "POST",
	}

	a2 := deepCopyRequestArgs(a)

	a2.method = "GET"

	assert.NotEqual(t, a.method, a2.method)
}

func TestDoPagination(t *testing.T) {
	// TODO: Mock out a local server to test rather than hitting github
	g := &ghAPI{
		baseURL: &url.URL{
			Host:   "api.github.com",
			Scheme: "https",
		},
		userAgent: "pr-test-code",
		version:   Version3,
		client:    &http.Client{},
	}

	a := &requestArgs{
		method:   "GET",
		endpoint: "/repos/coreos/etcd/pulls",
	}

	reqFunc := func(resp *http.Response) error {
		resp.Body.Close()
		return nil
	}

	err := g.doFullPagination(a, reqFunc)
	assert.NoError(t, err)
}
