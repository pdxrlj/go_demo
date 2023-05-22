package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cast"
)

func main() {

	http.HandleFunc("/tile", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		z := cast.ToInt(query.Get("z"))
		x := cast.ToInt(query.Get("x"))
		y := cast.ToInt(query.Get("y"))
		bytes, err := GetTileBytesDot3(z, y, x, "/Users/ruanlianjun/Desktop/0-3/_alllayers")
		if err != nil {
			panic(err)
		}
		if len(bytes) > 0 {
			contentType := http.DetectContentType(bytes[:50])
			writer.Header().Set("Content-Type", contentType)
			_, err = writer.Write(bytes)
			if err != nil {
				panic(err)
			}
			return
		}

		writer.WriteHeader(http.StatusNotFound)

	})

	log.Fatal(http.ListenAndServe(":8090", nil))

}

// GetTileBytesDot3 reads a tile from a local slice file for the corresponding level, row, and column,
// using the ArcGIS v2.0 slice data format.
func GetTileBytesDot3(mLevel, mRow, mColumn int, bundlesDir string) ([]byte, error) {
	const (
		tileSize   = 128
		headerSize = 64
	)
	var tileBytes []byte

	level := fmt.Sprintf("L%02d", mLevel)
	rowGroup := tileSize * (mRow / tileSize)
	row := fmt.Sprintf("R%04X", rowGroup)
	columnGroup := tileSize * (mColumn / tileSize)
	column := fmt.Sprintf("C%04X", columnGroup)
	bundleName := filepath.ToSlash(fmt.Sprintf("%s\\%s\\%s%s.bundle", bundlesDir, level, row, column))
	index := tileSize*(mRow-rowGroup) + (mColumn - columnGroup)
	fmt.Println("bundleName:", bundleName)
	if _, err := os.Stat(bundleName); os.IsNotExist(err) {
		return nil, nil
	}

	isBundle, err := os.Open(bundleName)
	if err != nil {
		return nil, err
	}
	defer isBundle.Close()
	isBundle.Seek(headerSize+8*int64(index), 0)
	indexBytes := make([]byte, 4)
	if _, err := isBundle.Read(indexBytes); err != nil {
		return nil, err
	}
	offset := int64(binary.LittleEndian.Uint32(indexBytes))
	startOffset := offset - 4
	lengthBytes := make([]byte, 4)
	if _, err := isBundle.ReadAt(lengthBytes, startOffset); err != nil {
		return nil, err
	}
	length := int(binary.LittleEndian.Uint32(lengthBytes))
	if length > 4 {
		tileBytes = make([]byte, length)
		if _, err := isBundle.Read(tileBytes); err != nil {
			return nil, err
		}
	}
	return tileBytes, nil
}
