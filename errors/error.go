package errors

import (
	"fmt"
)

type GPError struct {
	ErrMsg  string
	ErrCode int64
}

func New(message string, exts ...int64) *GPError {
	var code int64
	if len(exts) > 0 {
		code = exts[0]
	}
	return &GPError{
		ErrMsg:  message,
		ErrCode: code,
	}
}

func (p *GPError) Error() string {
	return p.ErrMsg
}

func (p *GPError) Code() int64 {
	return p.ErrCode
}

func (p *GPError) String() string {
	return fmt.Sprintf("code:%d message:%s", p.ErrCode, p.ErrMsg)
}
