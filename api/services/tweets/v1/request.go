package v1

import (
	"time"
)

// TweetsURLValues represents tweets URL query values.
type TweetsURLValues struct {
	AfterTime  time.Time `schema:"after_time"`
	BeforeTime time.Time `schema:"before_time"`
}
