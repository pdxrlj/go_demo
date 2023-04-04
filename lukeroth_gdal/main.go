package main

import (
	"fmt"
	"math"

	"github.com/lukeroth/gdal"
	geo "github.com/paulmach/go.geo"
)

func main() {
	dataset, err := gdal.Open("test.tif", gdal.ReadOnly)
	if err != nil {
		panic(err)
	}

	fmt.Println("dataset size: ", dataset.RasterXSize(), dataset.RasterYSize(), dataset.RasterCount())
	ori_transform := dataset.GeoTransform()

	XCount := dataset.RasterXSize()
	YCount := dataset.RasterYSize()

	latMin := ori_transform[3]
	latMax := ori_transform[3] + ori_transform[1]*float64(YCount)
	lonMin := ori_transform[0]
	lonMax := ori_transform[0] + ori_transform[1]*float64(XCount)
	fmt.Printf("latMin: %f, latMax: %f, lonMin: %f, lonMax: %f\n", latMin, latMax, lonMin, lonMax)
	spatialReference := gdal.CreateSpatialReference("")
	err = spatialReference.FromEPSGA(4326)
	if err != nil {
		panic(err)
	}
	imageBound := geo.NewBound(lonMax, lonMin, latMax, latMin)

	//4.2 获取原始影像的像素分辨率
	srcWEPixelResolution := (lonMax - lonMin) / float64(XCount)
	srcNSPixelResolution := (latMax - latMin) / float64(YCount)
	fmt.Printf("src_w_e_pixel_resolution: %f, src_n_s_pixel_resolution: %f\n", srcWEPixelResolution, srcNSPixelResolution)

	// 4.3 根据原始影像地理范围求解切片行列号  // 经纬度转瓦片编号
	zoom := 10
	tileRowMax := lat2tile(latMin, zoom) // 纬度  -90 ——90  lat
	tileRowMin := lat2tile(latMax, zoom) // 经度 -180 -- 180 lon
	tileColMin := lon2tile(lonMin, zoom)
	tileColMax := lon2tile(lonMax, zoom)
	fmt.Printf("tileRowMax: %d, tileRowMin: %d, tileColMin: %d, tileColMax: %d\n", tileRowMax, tileRowMin, tileColMin, tileColMax)
}

func lon2tile(lon float64, zoom int) int {
	return int(math.Floor((lon + 180) / 360 * math.Pow(2, float64(zoom))))
}

func lat2tile(lat float64, zoom int) int {
	return int(math.Floor((1 - math.Log(math.Tan(lat*math.Pi/180)+1/math.Cos(lat*math.Pi/180))/math.Pi) / 2 * math.Pow(2, float64(zoom))))
}

func tile2lon(col, zoom int) float64 {
	return float64(col)/math.Pow(2.0, float64(zoom))*360.0 - 180
}

func tile2lat(row, zoom int) float64 {
	n := math.Pi - (2.0*math.Pi*float64(row))/math.Pow(2.0, float64(zoom))
	return math.Atan(math.Sinh(n)) * 180 / math.Pi
}
