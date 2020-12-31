package json

import (
	"encoding/json"
	"io"
)

// NewDecoder returns a set up decoder.
func NewDecoder(body io.Reader) *json.Decoder {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	return decoder
}
