package api

import (
	"github.com/pkg/errors"
)

var (
	UserOrgErr = errors.New("Either User or Org must be set, not both")
)
