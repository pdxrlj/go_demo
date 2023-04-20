package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

var defaultHttpProxyUrl = "https://127.0.0.1:8080"

//Proxy-Connection

func main() {

	go server()

	proxyUrl, err := url.Parse(defaultHttpProxyUrl)
	if err != nil {
		panic(err)
	}
	defaultTransport := http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
		OnProxyConnectResponse: func(ctx context.Context, proxyURL *url.URL, connectReq *http.Request, connectRes *http.Response) error {
			return nil
		},
		GetProxyConnectHeader: func(ctx context.Context, proxyURL *url.URL, target string) (http.Header, error) {
			h := http.Header{
				"User-Agent": []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"},
				"Referer":    []string{"https://www.baidu.com"},
			}
			return h, nil
		},
		IdleConnTimeout: time.Second * 10,
	}
	client := http.Client{
		Transport: &defaultTransport,
	}

	response, err := client.Get("https://www.qq.com")
	if err != nil {
		panic(err)
	}
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	fmt.Printf("response:%v\n", string(bytes))

	select {}
}

func server() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		tmp := make([]byte, 0, 1024)
		n, err := conn.Read(tmp)
		if err != nil {
			panic(err)
		}
		fmt.Printf("read:%v\n", string(tmp[:n]))
	}
}
