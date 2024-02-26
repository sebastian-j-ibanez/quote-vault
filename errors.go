package main

import (
	"errors"
)

var (
	ErrMalformedQuery = errors.New("malformed query")
	ErrDbCreation     = errors.New("unable to initialize database")
	ErrOpenDb         = errors.New("unable to open database")
	ErrRowIteration   = errors.New("unable to iterate rows")
)
