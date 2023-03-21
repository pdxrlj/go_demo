package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

var videSrc = "https://shengubi-1304765822.cos.ap-nanjing.myqcloud.com/oceans.mp4"

//var videSrc = "https://shengubi-1304765822.cos.ap-nanjing.myqcloud.com/uploads/20210204/5d47836057e55606fb576453fb8c3237.jpg"

func main() {
	engine := gin.Default()

	engine.GET("/:file", func(context *gin.Context) {
		filenameBase := context.Param("file")
		filename := filepath.ToSlash(filenameBase)
		filename = fmt.Sprintf(videSrc, filename)

		httpReadSeeker := NewHttpReadSeeker(http.DefaultClient, videSrc)
		buf := make([]byte, 1024*1024*2)
		readSeeker := NewBufferedReadSeeker(httpReadSeeker, buf)
		http.ServeContent(context.Writer, context.Request, filenameBase, time.Now(), readSeeker)
	})

	engine.StaticFS("/html", gin.Dir(".", true))

	log.Fatalln(engine.Run(":8091"))
}
