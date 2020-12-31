package json

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Encoder struct {
	body   interface{}
	writer http.ResponseWriter
}

func NewEncoder(w http.ResponseWriter) *Encoder {
	return &Encoder{
		writer: w,
	}
}

func (enc *Encoder) SetBody(body interface{}) *Encoder {
	enc.body = body
	return enc
}

func (enc *Encoder) SetStatusInternalServerError() *Encoder {
	enc.writer.WriteHeader(http.StatusInternalServerError)
	return enc
}

func (enc *Encoder) SetStatusBadRequest() *Encoder {
	enc.writer.WriteHeader(http.StatusBadRequest)
	return enc
}

func (enc *Encoder) SetStatusNotFound() *Encoder {
	enc.writer.WriteHeader(http.StatusNotFound)
	return enc
}

func (enc *Encoder) SetStatusOK() *Encoder {
	enc.writer.WriteHeader(http.StatusOK)
	return enc
}

func (enc *Encoder) Encode() error {
	if enc.body == nil {
		return errors.New("missing body")
	}
	return json.NewEncoder(enc.writer).Encode(&enc.body)
}
