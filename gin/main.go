package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	engine.GET("/", func(context *gin.Context) {

	})

	log.Fatal(engine.Run(":8080"))
}
