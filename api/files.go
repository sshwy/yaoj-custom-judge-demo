package api

import (
	"errors"
	"log"

	"yaoj-go/service"
	"yaoj-go/utils"

	"github.com/gin-gonic/gin"
)

// 文件存储 REST API
type FileServiceRest struct {
	*service.FileService
}

var TempFiles = FileServiceRest{
	FileService: &service.TempFile,
}

var ContentLengthLimit = 5 * 1024 * 1024 // 5 MB

// Params:
// name string (query) 保存的文件名。如果是 random 则会随机生成名字
// ext string (query) 配合 name=random 情况使用，指定文件后缀名
// file []byte (body) 文件内容
// 返回一个可以 json 序列化的结构体
func (f *FileServiceRest) POST(ctx *gin.Context) (interface{}, error) {
	if ctx.Request.ContentLength > int64(ContentLengthLimit) {
		return nil, errors.New("content length limit exceed")
	}
	name := ctx.Query("name")
	if name == "" {
		return nil, errors.New("invalid name")
	} else if name == "random" {
		name = utils.RandString(12)
		if ext := ctx.Query("ext"); ext != "" {
			name = name + "." + ext
		}
	}

	if err := f.Store(name, ctx.Request.Body); err != nil {
		return nil, err
	}

	log.Printf("done name: %s\n", name)

	return gin.H{"path": name}, nil
}
