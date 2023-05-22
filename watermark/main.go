package main

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/golang/freetype"
)

//go:embed font.ttf
var fontStr string
var txt = "宏图科技"
var srcImage = "demo1.png"

func main() {
	source, err := os.Open(srcImage)
	if err != nil {
		panic(err)
	}
	srcImage, err := png.Decode(source)
	if err != nil {
		panic(err)
	}

	//1. 建立新的背景图层
	bgImage := image.NewRGBA(srcImage.Bounds())
	//2. 将原图像绘制到背景图层上
	for y := 0; y < bgImage.Bounds().Dy(); y++ {
		for x := 0; x < bgImage.Bounds().Dx(); x++ {
			bgImage.Set(x, y, srcImage.At(x, y))
		}
	}

	//3. 建立字体图层
	fontImage := image.NewRGBA(srcImage.Bounds())

	font, err := freetype.ParseFont([]byte(fontStr))
	if err != nil {
		panic(err)
	}
	f := freetype.NewContext()
	f.SetDPI(72)
	f.SetFont(font)
	f.SetFontSize(40)
	f.SetClip(fontImage.Bounds())
	f.SetDst(fontImage) // 设置写字的目标图层
	f.SetSrc(image.NewUniform(color.RGBA{R: 51, G: 200, A: 51}))
	pt := freetype.Pt(fontImage.Bounds().Max.X-40*4, fontImage.Bounds().Max.Y-5)
	fmt.Printf("X:%+v\n", pt)
	_, err = f.DrawString(txt, pt) // 将字写上去
	if err != nil {
		panic(err)
	}

	// 字体图层旋转30度
	//fontImg = imaging.Rotate(fontImg, 20, color.Transparent) // 右下角为支点旋转

	// 字体图层合并到背景图层
	draw.Draw(bgImage, bgImage.Bounds(), fontImage, image.Point{}, draw.Over)
	watermarkImg, err := os.Create("watermark.png")
	if err != nil {
		panic(err)
	}

	defer watermarkImg.Close()

	if err = png.Encode(watermarkImg, bgImage); err != nil {
		panic(err)
	}
}
