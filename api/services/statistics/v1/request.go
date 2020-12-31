package v1

import (
	"time"
)

// RelevanceURLValues represents relevance URL query values.
type RelevanceURLValues struct {
	AfterTime  time.Time `schema:"after_time"`
	BeforeTime time.Time `schema:"before_time"`
	Series     int       `schema:"series"`
}
