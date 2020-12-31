package v1

import "errors"

var (
	errInvalidSeries                = errors.New("invalid series")
	errInvalidBeforeTime            = errors.New("invalid before time")
	errInvalidBeforeTimeOrAfterTime = errors.New("invalid before time or after time")
	errInvalidAfterTime             = errors.New("invalid after time")
)
