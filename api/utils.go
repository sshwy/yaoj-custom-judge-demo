package api

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 可序列化的满足 error 接口的类型
type JsonError struct {
	err error
}

var _ error = (*JsonError)(nil)
var _ json.Marshaler = (*JsonError)(nil)

func (r JsonError) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(r.Error())), nil
}

func (r JsonError) Error() string {
	if r.err == nil {
		return ""
	}
	return r.err.Error()
}

func Jerr(s string) *JsonError {
	return &JsonError{err: errors.New(s)}
}

func GinHandler(f func(*gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := f(ctx)
		if err != nil {
			ctx.JSON(200, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(200, data)
		}
	}
}
