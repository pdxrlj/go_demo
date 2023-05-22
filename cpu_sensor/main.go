package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("sensors")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(output))
}
