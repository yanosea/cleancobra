package proxy

import (
	"encoding/json"
)

type Json interface {
	MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

type jsonProxy struct{}

func NewJson() Json {
	return &jsonProxy{}
}

func (jsonProxy) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func (jsonProxy) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
