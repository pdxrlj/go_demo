package main

import (
	"fmt"

	"github.com/deatil/go-cryptobin/cryptobin/crypto"
)

func main() {
	// 加密
	cypt := crypto.
		FromString("useData").
		SetKey("dfertf12dfertf12").
		Aes().
		ECB().
		PKCS7Padding().
		Encrypt().
		ToBase64String()

	// 解密
	cyptde := crypto.
		FromBase64String("i3FhtTp5v6aPJx0wTbarwg==").
		SetKey("dfertf12dfertf12").
		Aes().
		ECB().
		PKCS7Padding().
		Decrypt().
		ToString()

	// i3FhtTp5v6aPJx0wTbarwg==
	fmt.Println("加密结果：", cypt)
	fmt.Println("解密结果：", cyptde)

}
