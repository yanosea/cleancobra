package errors

import "fmt"

type wrappedError struct {
	msg string
	err error
}

func (e *wrappedError) Error() string {
	if e.err == nil && e.msg == "" {
		return ""
	} else if e.err == nil && e.msg != "" {
		return e.msg
	} else if e.err != nil && e.msg == "" {
		return e.err.Error()
	} else {
		return fmt.Sprintf("%s: %v", e.msg, e.err)
	}
}

func New(msg string) error {
	return &wrappedError{msg: msg}
}

func Wrap(err error, msg string) error {
	return &wrappedError{
		msg: msg,
		err: err,
	}
}
