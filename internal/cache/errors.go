package cache

import (
	"errors"
)

var (
	// ErrMiss means that value we are looking for wasn't in the cache.
	ErrMiss = errors.New("cache miss")
	// ErrLocalFailure means that there was an error while trying to access
	// the local cache.
	ErrLocalFailure = errors.New("error accessing local xkcd cache")
	// ErrOffline means that there was an error trying to access the xkcd
	// server.
	ErrOffline = errors.New("error accessing xkcd server")
)
