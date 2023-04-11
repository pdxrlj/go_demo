package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"

	"github.com/lukeroth/gdal"
)

var wg sync.WaitGroup

func Write2(bandCount, xOffset, yOffset, readWidth, readHeight, widthUsize, heightUsize int, path string, ds gdal.Dataset) {
	fmt.Println("img read success!")
	fmt.Println("band numbers", bandCount)
	fmt.Println("Image width", readWidth)
	fmt.Println("Image height", readHeight)

	img := make([][][]float64, bandCount) //储存影像数据的三维数组
	for i := 0; i < bandCount; i++ {
		img[i] = make([][]float64, readWidth)
	}

	for i := 0; i < bandCount; i++ {
		for j := 0; j < readWidth; j++ {
			img[i][j] = make([]float64, readHeight)
		}
	}

	minMax := make([][]float64, bandCount) //储存每个波段DN值的最大值与最小值

	wg.Add(bandCount)

	for i := 0; i < bandCount; i++ {
		minMax[i] = make([]float64, 2)
		band := ds.RasterBand(i + 1)
		go ReadDataFromBand(band, xOffset, yOffset, readWidth, readHeight, widthUsize, heightUsize, minMax[i], img[i])
	}

	wg.Wait()

	//程序运行到这里的时候影像数据已经读入到三维数组中了

	CreateImage(widthUsize, heightUsize, img, minMax, path) //将读入的数据输出为png

}

func ReadDataFromBand(band gdal.RasterBand, xOffset, yOffset, readWidth, readHeight, widthUsize, heightUsize int, minMax []float64, img [][]float64) {
	var tmp = make([]float64, widthUsize*heightUsize, widthUsize*heightUsize)

	err := band.IO(gdal.Read, xOffset, yOffset, readWidth, readHeight, tmp, widthUsize, heightUsize, 0, 0)
	if err != nil {
		panic(err)
	}

	minMax[0], _ = band.GetMinimum()
	minMax[1], _ = band.GetMaximum()
	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			if i+j*widthUsize > len(tmp) {
				img[i][j] = 0
				continue
			}
			img[i][j] = tmp[i+j*widthUsize]
		}
	}
	wg.Done()
}

func CreateImage(x, y int, imageArry [][][]float64, minMax [][]float64, path string) {
	width := x
	height := y

	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}

	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	// Set color for each pixel.
	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {

			/*
				R:=int(float64(imageArry[0][x][y]-minMax[0][0])/float64(minMax[0][1]-minMax[0][0])*float64(255))
				G:=int(float64(imageArry[1][x][y]-minMax[1][0])/float64(minMax[1][1]-minMax[1][0])*float64(255))
				B:=int(float64(imageArry[2][x][y]-minMax[2][0])/float64(minMax[2][1]-minMax[2][0])*float64(255))
				//线性拉伸
			*/

			color := color.RGBA{uint8(imageArry[0][x][y]), uint8(imageArry[1][x][y]), uint8(imageArry[2][x][y]), 0}
			img.Set(x, y, color)
		}
	}
	// Encode as PNG.
	f, _ := os.Create(path)
	png.Encode(f, img)
}
