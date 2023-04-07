package main

import (
	"fmt"
	"strings"
)

type Demo struct {
	Name string
	Age  int
}

func main() {
	//src := &Demo{
	//	Name: "pdx",
	//}
	//
	//dst := &Demo{
	//	Name: "rlj",
	//	Age:  2,
	//}
	//if err := mergo.Merge(src, dst); err != nil {
	//	panic(err)
	//}
	//
	//fmt.Printf("merge:%+v\n", src)

	path := "widgets_test/widgets.json"
	left := strings.Replace(path, "widgets_test/", "", 1)

	//if err != nil {
	//	panic(err)
	//}
	fmt.Println(left)

}
