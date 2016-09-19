package common

import (
	"fmt"
)

// ErrBadBody 不合理的请求内容.
type ErrBadBody struct {
	Message string
}

func (e ErrBadBody) Error() string { return "请求非法-" + e.Message }

// BadBodyErrf 非法请求内容.
func BadBodyErrf(format string, a ...interface{}) ErrBadBody {
	return ErrBadBody{Message: fmt.Sprintf(format, a...)}
}

// BadBodyErr 非法请求内容.
func BadBodyErr(a ...interface{}) ErrBadBody {
	return ErrBadBody{Message: fmt.Sprint(a...)}
}
