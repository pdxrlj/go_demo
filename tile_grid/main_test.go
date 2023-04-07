package main

import (
	"fmt"
	"math"
	"testing"

	"github.com/lukeroth/gdal"
	"github.com/magiconair/properties/assert"
)

func TestResolution(t *testing.T) {
	resolution := Resolution(10)
	assert.Equal(t, resolution, 152.8740565703525)
}
func TestResolution2(t *testing.T) {
	initialResolution := 2 * math.Pi * 6378137 / 256
	var resolution = initialResolution / math.Pow(2, 10)
	fmt.Printf("resolution: %v\n", resolution) //152.8740565703525
}

func TestLonLatToMercator(t *testing.T) {
	x, y := LonLatToMercator(116.3, 39.85)
	fmt.Printf("x: %v, y: %v\n", x, y)
	assert.Equal(t, x, 12946456.779257717)
	assert.Equal(t, y, 4844168.570470276)
}

func TestComputeStartXY(t *testing.T) {
	x, y := LonLatToMercator(120.75141906738281, 30.759538817987497)

	// {level:16 levelTileMax:65535 minX:54750 minY:26878 maxX:54770 maxY:26892}
	x, y = ComputeStartXY(x, y, 16)
	fmt.Printf("x: %v, y: %v\n", x, y)
}

func TestMinMaxZoom(t *testing.T) {
	dataset, err := gdal.Open("/Users/ruanlianjun/Desktop/demo.tif", gdal.ReadOnly)
	if err != nil {
		panic(err)
	}
	transform := dataset.GeoTransform()
	fmt.Printf("geotransform: %v\n", transform)

	minZoom, maxZoom := MinMaxZoom(113.1710635, 114.4020218, float64(dataset.RasterXSize()), float64(dataset.RasterXSize()))
	fmt.Printf("minZoom: %v, maxZoom: %v\n", minZoom, maxZoom)
}

func TestComputeStartGdalXY(t *testing.T) {
	dataset, err := gdal.Open("/Users/ruanlianjun/Desktop/demo.tif", gdal.ReadOnly)
	if err != nil {
		panic(err)
	}

	transform := dataset.GeoTransform()

	x, y := ComputeStartXY(transform[0], transform[2], 2)
	fmt.Printf("x: %v, y: %v\n", x, y)
}
