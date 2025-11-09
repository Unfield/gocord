package gateway

import "errors"

var ErrNoSession = errors.New("no active session to resume")
