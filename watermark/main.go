package main

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
)

//go:embed font2.ttf
var fontStr string
var txt = "宏图科技"
var srcImage = "bg1.png"

var fontImage *image.RGBA

func init() {
	//3. 建立字体图层
	fontImage = image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: 256, Y: 256},
	})

	font, err := freetype.ParseFont([]byte(fontStr))
	if err != nil {
		panic(err)
	}
	fmt.Printf("fontImage.Bounds():%+v\n", fontImage.Bounds())
	f := freetype.NewContext()
	//f.SetDPI(20)
	f.SetFont(font)
	f.SetFontSize(30)
	f.SetClip(fontImage.Bounds())
	f.SetDst(fontImage) // 设置写字的目标图层
	f.SetSrc(image.NewUniform(color.RGBA{R: 255, G: 0, B: 0, A: 255}))
	pt := freetype.Pt((fontImage.Bounds().Max.X)/2, (fontImage.Bounds().Max.Y)/2)
	_, err = f.DrawString(txt, pt)
	if err != nil {
		panic(err)
	}

	// 字体图层旋转30度
	fontImage = (*image.RGBA)(imaging.Rotate(fontImage, 50, color.Transparent))
}

func main() {
	source, err := os.Open(srcImage)
	if err != nil {
		panic(err)
	}
	srcImageDecode, err := png.Decode(source)
	if err != nil {
		panic(err)
	}

	//1. 建立新的背景图层
	bgImage := image.NewRGBA(srcImageDecode.Bounds())
	//2. 将原图像绘制到背景图层上
	for y := 0; y < bgImage.Bounds().Dy(); y++ {
		for x := 0; x < bgImage.Bounds().Dx(); x++ {
			bgImage.Set(x, y, srcImageDecode.At(x, y))
		}
	}

	// 字体图层合并到背景图层

	draw.Draw(bgImage, srcImageDecode.Bounds(), fontImage, image.Point{}, draw.Over)

	watermarkImg, err := os.Create("watermark.png")
	if err != nil {
		panic(err)
	}
	defer watermarkImg.Close()

	if err = png.Encode(watermarkImg, bgImage); err != nil {
		panic(err)
	}
}
