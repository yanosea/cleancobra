//go:generate mockgen -source=strconv.go -destination=strconv_mock.go -package=proxy

package proxy

import (
	"strconv"
)

// Strconv provides a proxy interface for strconv package functions
type Strconv interface {
	Atoi(s string) (int, error)
	Itoa(i int) string
	ParseBool(str string) (bool, error)
	ParseInt(s string, base int, bitSize int) (int64, error)
	ParseFloat(s string, bitSize int) (float64, error)
	FormatBool(b bool) string
	FormatInt(i int64, base int) string
	FormatFloat(f float64, fmt byte, prec, bitSize int) string
}

// StrconvImpl implements the Strconv interface using the standard strconv package
type StrconvImpl struct{}

// NewStrconv creates a new Strconv implementation
func NewStrconv() Strconv {
	return &StrconvImpl{}
}

func (s *StrconvImpl) Atoi(str string) (int, error) {
	return strconv.Atoi(str)
}

func (s *StrconvImpl) Itoa(i int) string {
	return strconv.Itoa(i)
}

func (s *StrconvImpl) ParseBool(str string) (bool, error) {
	return strconv.ParseBool(str)
}

func (s *StrconvImpl) ParseInt(str string, base int, bitSize int) (int64, error) {
	return strconv.ParseInt(str, base, bitSize)
}

func (s *StrconvImpl) ParseFloat(str string, bitSize int) (float64, error) {
	return strconv.ParseFloat(str, bitSize)
}

func (s *StrconvImpl) FormatBool(b bool) string {
	return strconv.FormatBool(b)
}

func (s *StrconvImpl) FormatInt(i int64, base int) string {
	return strconv.FormatInt(i, base)
}

func (s *StrconvImpl) FormatFloat(f float64, fmt byte, prec, bitSize int) string {
	return strconv.FormatFloat(f, fmt, prec, bitSize)
}