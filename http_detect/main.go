package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	content, err := os.ReadFile("47.octet-stream")
	if err != nil {
		panic(err)
	}
	contentType := http.DetectContentType(content[:50])
	fmt.Println(contentType)
}
