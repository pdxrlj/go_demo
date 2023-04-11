package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

func WriteFinal(width, height int, imageData []byte, filename string) {
	fmt.Println("width:", width, "height:", height)
	img := image.NewRGBA(image.Rect(0, 0, 256, 256))

	var pixels [][][]uint8
	for i := 0; i < height; i++ {
		row := make([][]uint8, width)
		for j := 0; j < width; j++ {
			pixel := make([]uint8, 3)
			for k := 0; k < 3; k++ {
				if (i*width+j)*3+k >= len(imageData) {
					pixel[k] = 0
					continue
				}
				pixel[k] = imageData[(i*width+j)*3+k]
			}
			row[j] = pixel
		}
		pixels = append(pixels, row)
	}

	// 将 3D 数组转换为 Go 语言中的 image.RGBA 类型
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			pixel := color.RGBA{
				R: pixels[i][j][0],
				G: pixels[i][j][1],
				B: pixels[i][j][2],
				A: 255,
			}
			img.SetRGBA(j, i, pixel)
		}
	}

	// 将图像编码为PNG格式，并将结果写入文件。
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}

}

func Write(bandCount, width, height int, data []byte, fileName string) {
	fmt.Printf("data len: %d\n", len(data))
	fmt.Printf("data width len: %d\n", width)
	fmt.Printf("data height len: %d\n", height)
	img := make([][][]byte, bandCount) //储存影像数据的三维数组
	for i := 0; i < bandCount; i++ {
		img[i] = make([][]byte, 256)
	}

	for i := 0; i < bandCount; i++ {
		for j := 0; j < 256; j++ {
			img[i][j] = make([]byte, height)
		}
	}

	for i := 0; i < bandCount; i++ {
		ReadData(width, height, data, img[i])
	}

	writeImage(width, height, img, fileName)
}

func writeImage(width, height int, imageArray [][][]byte, fileName string) {
	upLeft := image.Point{}
	lowRight := image.Point{X: 256, Y: 256}

	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	// Set color for each pixel.
	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			if x >= len(imageArray[0]) || x >= len(imageArray[1]) || x >= len(imageArray[2]) {
				img.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 0xff})
				continue
			}
			if y >= len(imageArray[0][x]) || y >= len(imageArray[1][x]) || y >= len(imageArray[2][x]) {
				img.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 0xff})
				continue
			}

			rgba := color.RGBA{R: imageArray[0][x][y], G: imageArray[1][x][y], B: imageArray[2][x][y], A: 150}
			img.Set(x, y, rgba)
		}
	}
	// Encode as PNG.
	f, _ := os.Create(fileName)
	err := png.Encode(f, img)
	if err != nil {
		return
	}
}

func ReadData(width, height int, data []byte, imageArray [][]byte) {
	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			imageArray[i][j] = data[i+j*256]
		}
	}
}
