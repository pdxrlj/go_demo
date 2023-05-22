package main

import (
	"fmt"
	"path/filepath"
)

func main() {

	fmt.Println(filepath.Base("a/v.txt"))

	//	fmt.Println("start")
	//	goto DONE
	//	fmt.Println("end")
	//	goto DONE
	//
	//DONE:
	//	fmt.Println("done")

	//file, err := os.Open("go.mod")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//tmp := make([]byte, 10)
	//firstLine, err := file.Read(tmp)
	//fmt.Println("1:", string((tmp[:firstLine])))
	//
	//tmp2 := make([]byte, 1000)
	//secondLine, err := file.Read(tmp2)
	//fmt.Println("2:", string(tmp2[:secondLine]))

	//for i := 0; i < 10; i++ {
	//	for j := 0; j < 30; j++ {
	//		if j == 5 {
	//			break
	//		}
	//		fmt.Println("i:", i, "j:", j)
	//	}
	//}
}
