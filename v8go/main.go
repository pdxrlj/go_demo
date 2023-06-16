package main

import (
	"fmt"
	"time"

	"rogchap.com/v8go"
)

func main() {

	resp := make(chan struct{})

	go func() {
		time.Sleep(time.Second * 2)
		resp <- struct{}{}
	}()

	select {
	case <-time.After(time.Second * 5):
		fmt.Println("timeout")
	case <-resp:
		fmt.Println("default")
	}

	return

	// 设置一个 JavaScript 函数，该函数接收两个参数并将它们相加
	jsAddFunc := `
		function add(a, b) {
			return a + b;
		}
	`
	ctx := v8go.NewContext()
	iso := ctx.Isolate()
	defer iso.Dispose()
	defer ctx.Close()

	_, err := ctx.RunScript(jsAddFunc, "")
	if err != nil {
		panic(err)
	}
	value, err := ctx.Global().Get("add")
	if err != nil {
		panic(err)
	}

	v1, err := v8go.NewValue(iso, int32(1))
	if err != nil {
		panic(err)
	}
	v2, err := v8go.NewValue(iso, int32(2))
	if err != nil {
		panic(err)
	}

	// 调用 JavaScript 函数
	result, err := value.AsFunction()
	if err != nil {
		panic(err)
	}
	call, err := result.Call(v8go.Undefined(iso), v1, v2)
	if err != nil {
		panic(err)
	}
	fmt.Println(call.String())
}
