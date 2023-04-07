package main

import (
	"archive/zip"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("/Users/ruanlianjun/Desktop/widgets_test.zip")
	if err != nil {
		panic(err)
	}
	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	reader, err := zip.NewReader(file, stat.Size())
	if err != nil {
		panic(err)
	}
	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		fmt.Printf("file name: %s\n", file.Name)
	}
}
