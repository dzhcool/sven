package errors

var (
	// system
	ErrNil          = &GPError{ErrCode: 1, ErrMsg: "nil or type error"}
	ErrConfEmpty    = &GPError{ErrCode: 1, ErrMsg: "conf empty"}
	ErrFileNotExist = &GPError{ErrCode: 1, ErrMsg: "file not exist"}
	ErrExpired      = &GPError{ErrCode: 1, ErrMsg: "expired"}

	ErrSearchFailed = &GPError{ErrCode: 1, ErrMsg: "查询数据失败"}
	ErrUpdateFailed = &GPError{ErrCode: 1, ErrMsg: "更新数据失败"}
)
