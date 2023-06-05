package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func main() {
	x, y := robotgo.GetScreenSize()
	fmt.Println("color----", x, y)
}
