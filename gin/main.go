package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	url := "https://www.baidu.com"

	before, _, found := strings.Cut(url, "?")
	fmt.Println("before:", before, "found:", found)
	return

	engine := gin.Default()
	engine.GET("/:id/:name", func(context *gin.Context) {
		id := context.Param("id")
		name := context.Param("name")
		fmt.Println("id:", id, "name:", name)
		sex := context.Query("sex")
		context.String(200, "id:"+id+" name:"+name+" sex:"+sex)
	})

	engine.GET("/index.html", func(context *gin.Context) {
		file, err := os.Open("index.html")
		if err != nil {
			context.String(404, "%s", err.Error())
			return
		}
		all, err := io.ReadAll(file)
		if err != nil {
			context.String(404, "%s", err.Error())
			return
		}
		context.String(200, "%s", all)
	})

	log.Fatal(engine.Run(":18085"))
}
