package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/get", func(context echo.Context) error {
		fmt.Printf("get query:%+v\n", context.QueryParams())
		fmt.Printf("get path:%+v\n", context.Path())
		return context.String(http.StatusOK, "get")
	})

	e.POST("/post", func(context echo.Context) error {
		fmt.Printf("post query:%+v\n", context.QueryParams())
		fmt.Printf("post path:%+v\n", context.Path())
		values, err := context.FormParams()
		if err != nil {
			fmt.Printf("post form:%+v\n", err)
			return err
		}
		fmt.Printf("post form:%+v\n", values)
		return context.String(http.StatusOK, "post")
	})
	log.Fatal(e.Start(":8083"))
}
