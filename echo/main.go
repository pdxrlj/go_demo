package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

func main() {
	u := "http://117.78.4.140:3389/tile-server/v1?layer=China_08m_DOM_3857_WMTS^&tilematrix={z}^&Tilecol={x}^&TileRow={y}"
	escapeUrl, err := url.QueryUnescape(u)
	fmt.Printf("escapeUrl:%+v err:%v\n", escapeUrl, err)
	return

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
