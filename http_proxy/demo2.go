package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

func main() {
	// tcp 连接，监听 8080 端口
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Panic(err)
	}
	// 死循环，每当遇到连接时，调用 handle
	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handle2(client)
	}
}

var proxy5 = socket5()

func socket5() proxy.Dialer {
	socks5, err := proxy.SOCKS5("tcp", "127.0.0.1:7890", nil, proxy.Direct)
	if err != nil {
		panic(err)
	}

	return socks5
}

func httpClient(servers net.Conn) *http.Client {
	c := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return servers, nil
			},
			MaxIdleConnsPerHost: 1000,
			MaxIdleConns:        1000,
		},
		Timeout: time.Second * 10,
	}
	return c
}

var schema = "http://"

func handle2(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()
	log.Printf("remote addr: %v\n", client.RemoteAddr())
	// 用来存放客户端数据的缓冲区
	var b = make([]byte, 0, 1024)
	var err error
	for {
		var tmp [1024]byte
		n, err := client.Read(tmp[:])

		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Printf("error:%T, err:%v\n", err, err)
				return
			}
		}
		if n == 0 {
			break
		}

		b = append(b, tmp[:n]...)
		if n < len(tmp) {
			break
		}
	}
	var method, address string

	request, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(b)))
	if err != nil {
		log.Println(err)
		return
	}
	method = request.Method
	schema = request.URL.Scheme
	host := request.URL.Host

	if err != nil {
		log.Println(err)
		return
	}
	if method == "CONNECT" {
		address = host
	} else {
		address = host
		if strings.Index(host, ":") == -1 {
			address = host + ":80"
		}
	}

	servers, err := proxy5.Dial("tcp", address)
	if method == "CONNECT" {
		_, _ = fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		_, _ = servers.Write(b)

		u, _ := url.Parse(address)
		r := &http.Request{
			Method: method,
			URL:    u,
		}

		resp, err := httpClient(servers).Do(r)
		if err != nil {
			log.Println("do:", err)
			return
		}
		fmt.Printf("resp:%v\n", resp.Trailer)
	}
	go io.Copy(servers, client)
	if servers != nil {
		io.Copy(client, servers)
	}
}
