package verr

import (
	"fmt"
	"io"
	"strings"
)

var StackEnabled = true

// Error severance level type. Can be used to determine the importance of the error before it handled appropriately.
type ErrLevel int

const (
	ErrError ErrLevel = iota
	ErrIgnorable
	ErrInfo
	ErrNotice
	ErrWarning
	ErrCritical
	ErrFatal
	ErrPanic
)

type VError struct {
	Code      int
	Level     ErrLevel
	Msg       string
	Cause     error
	Stack     *stack
	ErrParams map[string]interface{}
}

var LevelMap = map[ErrLevel]string{
	ErrIgnorable: "ignorable",
	ErrInfo:      "info",
	ErrNotice:    "notice",
	ErrWarning:   "warning",
	ErrError:     "error",
	ErrCritical:  "critical",
	ErrFatal:     "fatal",
	ErrPanic:     "panic",
}

func (fe *VError) Error() string {
	text := []string{fe.Msg}
	curErr := fe.Cause
	for curErr != nil {
		if e, ok := fe.Cause.(*VError); ok {
			text = append(text, e.Msg)
			curErr = e.Cause
		} else {
			text = append(text, curErr.Error())
			break
		}
	}
	return strings.Join(text, ": ")
}

func (fe *VError) WithLevel(level ErrLevel) *VError {
	fe.Level = level
	return fe
}

func (fe *VError) WithCause(cause error) *VError {
	fe.Cause = cause
	return fe
}

func (fe *VError) WithCode(code int) *VError {
	fe.Code = code
	return fe
}

func (fe *VError) AddParam(key string, v interface{}) *VError {
	if fe.ErrParams == nil {
		fe.ErrParams = make(map[string]interface{})
	}
	fe.ErrParams[key] = v
	return fe
}

func (fe *VError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", fe.Cause)
			fe.Stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, fe.Error())
	case 'q':
		fmt.Fprintf(s, "%q", fe.Error())
	}
}

func Error(message string) *VError {
	return &VError{
		Level: ErrError,
		Msg:   message,
		Stack: callers(),
	}
}

func Errorf(format string, args ...interface{}) *VError {
	return &VError{
		Level: ErrError,
		Msg:   fmt.Sprintf(format, args...),
		Stack: callers(),
	}
}

func Wrap(err error, message string) *VError {
	return &VError{
		Msg:   message,
		Stack: callers(),
		Cause: err,
	}
}

func Cause(err error) error {
	e, ok := err.(*VError)
	if ok {
		return e.Cause
	}
	return nil
}

func Params(err error) map[string]interface{} {
	e, ok := err.(*VError)
	if ok {
		return e.ErrParams
	}
	return nil
}

func Level(err error) ErrLevel {
	e, ok := err.(*VError)
	if ok {
		return e.Level
	}
	return ErrError
}

func Code(err error) int {
	e, ok := err.(*VError)
	if ok {
		return e.Code
	}
	return -1
}
