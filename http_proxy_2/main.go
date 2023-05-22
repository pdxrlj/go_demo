package main

import (
	"http_proxy_v2/proxy"
)

func main() {
	proxy.Serve(":8083")
}
