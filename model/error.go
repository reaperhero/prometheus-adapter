package model

import "errors"

const ErrCodeOK = 1000

var (
	ErrInvalidParam = errors.New("参数不合法")
	ErrDbOperation  = errors.New("数据库错误")
)

func GetErrorCode(err error) int32 {
	switch err {
	case ErrInvalidParam:
		return ErrCodeOK + 1
	case ErrDbOperation:
		return ErrCodeOK + 2
	case nil:
		return ErrCodeOK
	default:
		return ErrCodeOK
	}
}

func GetErrorMap(err error) map[string]interface{} {
	var msg = "OK"
	if err != nil {
		msg = err.Error()
	}

	return map[string]interface{}{
		"code": GetErrorCode(err),
		"msg":  msg,
	}
}
