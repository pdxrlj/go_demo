package main

import (
	"fmt"

	"github.com/robertkrimen/otto"
)

const SCRIPT = `
function main(param) {
	obj= {"name":"demo", "age":18, "param":param, "arr":["1", "2", "3"], "obj":{"name":"demo"}};
	return obj;
};
main("123")
`

type MyStruct struct {
	Val  int
	Name string
	Age  int
}

func main() {

	vm := otto.New()
	value, err := vm.Run(SCRIPT)

	if err != nil {
		panic(err)
	}
	fmt.Printf("value: %+v\n", value)

	//demo := &MyStruct{
	//	Val:  1,
	//	Name: "demo",
	//	Age:  18,
	//}
	//vm := goja.New()
	//script, err := vm.RunString(SCRIPT)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Printf("script: %+v script: %v\n", script, script.Export())

}
