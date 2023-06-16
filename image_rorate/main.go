package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/pdxrlj/graphics"
)

func ParseHexColor(fontColor string) (color.RGBA, error) {
	return hex2rgb(fontColor), nil
}

func hex2rgb(str string) color.RGBA {
	str = strings.Trim(str, "#")
	r, _ := strconv.ParseInt(str[:2], 16, 10)
	g, _ := strconv.ParseInt(str[2:4], 16, 10)
	b, _ := strconv.ParseInt(str[4:], 16, 10)
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
}

func parseHexColor(hexColor string) (color.RGBA, error) {
	// 去除 "#" 前缀
	hex := hexColor[1:]

	// 解析 HexColor
	rgb, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return color.RGBA{}, err
	}

	// 将解析结果转换为 RGBA
	rgba := color.RGBA{
		R: uint8(rgb >> 16),
		G: uint8(rgb >> 8),
		B: uint8(rgb),
		A: 255,
	}

	return rgba, nil
}

func RotateImage(img image.Image, angle int) (*image.RGBA, error) {
	angle = angle % 360
	// 弧度转换
	radian := float64(angle) * math.Pi / 180.0
	cos := math.Cos(radian)
	sin := math.Sin(radian)
	// 原图的宽高
	w := float64(img.Bounds().Dx())
	h := float64(img.Bounds().Dy())

	// 新图高宽
	W := int(math.Max(math.Abs(w*cos-h*sin), math.Abs(w*cos+h*sin)))
	H := int(math.Max(math.Abs(w*sin-h*cos), math.Abs(w*sin+h*cos)))

	dst := image.NewRGBA(image.Rect(0, 0, W, H))
	if err := graphics.Rotate(dst, img, &graphics.RotateOptions{Angle: radian}); err != nil {
		return nil, err
	}

	return dst, nil
}

func main() {
	open, err := os.Open("test.png")
	if err != nil {
		panic(err)
	}
	decode, _, err := image.Decode(open)
	if err != nil {
		panic(err)
	}

	rotateImage, err := RotateImage(decode, 30)
	if err != nil {
		panic(err)
	}

	out, err := os.Create("out.png")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	err = png.Encode(out, rotateImage)
	if err != nil {
		panic(err)
	}
}
