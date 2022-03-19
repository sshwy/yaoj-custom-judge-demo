package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"yaoj-go/api"
	"yaoj-go/service"

	"github.com/gin-gonic/gin"
	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "welcome")
}

func AllowAller() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")

		// before request

		c.Next()

		// after request
	}
}

func main() {
	// register rpc service
	if err := rpc.Register(&service.Judge); err != nil {
		log.Fatal(err)
	}

	// set up server
	r := gin.Default()
	r.Use(AllowAller())
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "welcome")
	})
	r.POST("/files", api.GinHandler(api.TempFiles.POST))
	r.POST("/jsonrpc", func(ctx *gin.Context) {
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			ReadCloser: ctx.Request.Body,
			Writer:     ctx.Writer,
		}

		ctx.Header("Content-Type", "application/json")

		if err := rpc.ServeRequest(jsonrpc.NewServerCodec(conn)); err != nil {
			log.Println(err)
		}
	})
	// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Methods/OPTIONS
	r.OPTIONS("/jsonrpc", func(ctx *gin.Context) {
		ctx.Status(200)
		ctx.Header("Allow", "POST")
		ctx.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "X-PINGOTHER, Content-Type")
	})

	r.Run(":3000")
}
