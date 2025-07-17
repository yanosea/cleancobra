//go:generate mockgen -source=json.go -destination=json_mock.go -package=proxy

package proxy

import (
	"encoding/json"
	"io"
)

// JSON provides a proxy interface for encoding/json package functions
type JSON interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
	NewEncoder(w io.Writer) JSONEncoder
	NewDecoder(r io.Reader) JSONDecoder
}

// JSONEncoder provides a proxy interface for json.Encoder
type JSONEncoder interface {
	Encode(v any) error
	SetIndent(prefix, indent string)
}

// JSONDecoder provides a proxy interface for json.Decoder
type JSONDecoder interface {
	Decode(v any) error
}

// JSONImpl implements the JSON interface using the standard encoding/json package
type JSONImpl struct{}

// JSONEncoderImpl implements the JSONEncoder interface
type JSONEncoderImpl struct {
	encoder *json.Encoder
}

// JSONDecoderImpl implements the JSONDecoder interface
type JSONDecoderImpl struct {
	decoder *json.Decoder
}

// NewJSON creates a new JSON implementation
func NewJSON() JSON {
	return &JSONImpl{}
}

func (j *JSONImpl) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j *JSONImpl) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (j *JSONImpl) NewEncoder(w io.Writer) JSONEncoder {
	return &JSONEncoderImpl{encoder: json.NewEncoder(w)}
}

func (j *JSONImpl) NewDecoder(r io.Reader) JSONDecoder {
	return &JSONDecoderImpl{decoder: json.NewDecoder(r)}
}

func (e *JSONEncoderImpl) Encode(v any) error {
	return e.encoder.Encode(v)
}

func (e *JSONEncoderImpl) SetIndent(prefix, indent string) {
	e.encoder.SetIndent(prefix, indent)
}

func (d *JSONDecoderImpl) Decode(v any) error {
	return d.decoder.Decode(v)
}