package proxy

import (
	"encoding/json"
)

type Json interface {
	MarshalIndent(v any, prefix, indent string) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

type jsonProxy struct{}

func NewJson() Json {
	return &jsonProxy{}
}

func (jsonProxy) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func (jsonProxy) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
