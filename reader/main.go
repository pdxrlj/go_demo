package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("go.mod")
	if err != nil {
		panic(err)
	}

	err = BufIo(file)
	if err != nil {
		panic(err)
	}
}

func BufIo(read io.Reader) error {
	scanner := bufio.NewScanner(read)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		fmt.Printf("%s\n", string(bytes))
		fmt.Println(strings.Repeat("=", 10))
	}

	return scanner.Err()
}
