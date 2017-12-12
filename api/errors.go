package api

import (
	"github.com/pkg/errors"
)

var (
	// ErrUserOrg reports if the user and/or organization specified is not valid
	ErrUserOrg = errors.New("Either User or Org must be set, not both")
)
