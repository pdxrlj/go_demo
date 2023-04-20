package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "hello world")
		fmt.Println("=========", request.Method)
	})
	http.ListenAndServe(":8089", nil)
}
