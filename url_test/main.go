package main

import (
	"fmt"
	"net/url"
)

const (
	urlTest = "amqp://guest:guest@localhost:5672/"
)

func main() {
	parse, err := url.Parse(urlTest)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v", parse)
}
