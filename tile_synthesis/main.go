package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/lukeroth/gdal"
	"github.com/pkg/errors"
)

func main() {
	err := NewBuildGTiff(
		WithBuildGTiffBounds(6761, 3329, 6774, 3340),
		WithBuildGTiffDestFilename("build.tiff"),
		WithBuildGTiffLevel(13),
		WithBuildGTiffTileSize(256),
		WithBuildGTiffBand(3),
		WithBuildGTiffSourceDir("/Users/ruanlianjun/Desktop/tile_download_test/13"),
	).Init().Build()
	if err != nil {
		fmt.Printf("err: %+v\n", err)
	}
}

type BuildGTiff struct {
	MinX, MinY, MaxX, MaxY int
	TileSize               int
	Level                  int
	DestFilename           string
	Band                   int
	sourceDir              string

	outDs *gdal.Dataset
}

type BuildGTiffOption func(*BuildGTiff)

func WithBuildGTiffBounds(minx, miny, maxx, maxy int) BuildGTiffOption {
	return func(build *BuildGTiff) {
		build.MinX = minx
		build.MinY = miny
		build.MaxX = maxx
		build.MaxY = maxy
	}
}

func WithBuildGTiffTileSize(tileSize int) BuildGTiffOption {
	return func(build *BuildGTiff) {
		build.TileSize = tileSize
	}
}

func WithBuildGTiffLevel(level int) BuildGTiffOption {
	return func(build *BuildGTiff) {
		build.Level = level
	}
}

func WithBuildGTiffDestFilename(destFilename string) BuildGTiffOption {
	return func(build *BuildGTiff) {
		build.DestFilename = destFilename
	}
}

func WithBuildGTiffBand(band int) BuildGTiffOption {
	return func(build *BuildGTiff) {
		build.Band = band
	}
}

func WithBuildGTiffSourceDir(sourceDir string) BuildGTiffOption {
	return func(build *BuildGTiff) {
		build.sourceDir = sourceDir
	}
}

func NewBuildGTiff(options ...BuildGTiffOption) *BuildGTiff {
	b := &BuildGTiff{}
	for _, option := range options {
		option(b)
	}
	return b
}

func (b *BuildGTiff) Init() *BuildGTiff {
	driver, err := gdal.GetDriverByName("GTiff")
	if err != nil {
		panic(err)
	}
	geoTransForm := NewTileCoordinateBound(
		WithBounds(b.MinX, b.MinY, b.MaxX, b.MaxY),
		WithTileSize(b.TileSize),
		WithLevel(b.Level),
	).GetGeoTransform()

	resultWidth := int(math.Ceil(float64(b.MaxX-b.MinX+1) * float64(b.TileSize)))
	resultHeight := int(math.Ceil(float64(b.MaxY-b.MinY+1) * float64(b.TileSize)))
	fmt.Printf("resultWidth: %d, resultHeight: %d\n", resultWidth, resultHeight)
	_ = os.Remove(b.DestFilename)
	create := driver.Create(b.DestFilename, resultWidth, resultHeight, b.Band, gdal.Byte, nil)
	// 设置crate的geometry
	err = create.SetGeoTransform(geoTransForm)
	if err != nil {
		panic("SetGeoTransform: " + err.Error())
	}
	b.outDs = &create
	return b
}

func (b *BuildGTiff) Build() error {
	defer b.outDs.Close()
	err := filepath.WalkDir(b.sourceDir, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		y := filepath.Base(path)[:strings.Index(filepath.Base(path), ".")]
		x := filepath.Base(filepath.Dir(path))
		//z := filepath.Base(filepath.Dir(filepath.Dir(path)))
		fmt.Printf("x: %s, y: %s\n", x, y)
		xInt, err := strconv.Atoi(x)
		if err != nil {
			return err
		}

		yInt, err := strconv.Atoi(y)
		if err != nil {
			return err
		}
		xOffset := (xInt - b.MinX) * b.TileSize
		yOffset := (yInt - b.MinY) * b.TileSize

		dataset, err := gdal.Open(path, gdal.ReadOnly)
		if err != nil {
			return err
		}
		defer dataset.Close()
		band := dataset.RasterCount()
		width := dataset.RasterXSize()
		height := dataset.RasterYSize()
		for i := 0; i < band; i++ {
			data := make([]byte, width*height)
			rasterBand := dataset.RasterBand(i + 1)
			err := rasterBand.IO(gdal.Read, 0, 0, width, height, data, width, height, 0, 0)
			if err != nil {
				return errors.WithStack(err)
			}

			err = b.outDs.RasterBand(i+1).IO(gdal.Write, xOffset, yOffset, width, height, data, width, height, 0, 0)
			if err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})

	b.outDs.FlushCache()

	return err
}

type TileCoordinateBound struct {
	MinX, MinY, MaxX, MaxY int
	initialResolution      float64
	tileSize               int
	zoom                   int
	// 地球周长的一半
	originShift float64
	resolution  float64

	rMinX, rMinY, rMaxX, rMaxY float64
}

type TileCoordinateOption func(*TileCoordinateBound)

func WithTileSize(tileSize int) TileCoordinateOption {
	return func(t *TileCoordinateBound) {
		t.tileSize = tileSize
	}
}

func WithLevel(level int) TileCoordinateOption {
	return func(t *TileCoordinateBound) {
		t.zoom = level
	}
}

func WithBounds(minx, miny, maxx, maxy int) TileCoordinateOption {
	return func(t *TileCoordinateBound) {
		t.MinX = minx
		t.MinY = miny
		t.MaxX = maxx
		t.MaxY = maxy
	}
}

func NewTileCoordinateBound(options ...TileCoordinateOption) *TileCoordinateBound {
	t := &TileCoordinateBound{
		tileSize: 256,
	}
	t.initialResolution = 2 * math.Pi * 6378137 / float64(t.tileSize)
	t.originShift = 2 * math.Pi * 6378137 / 2.0
	for _, option := range options {
		option(t)
	}

	return t
}

func (t *TileCoordinateBound) Resolution(z int) *TileCoordinateBound {
	t.resolution = t.initialResolution / math.Pow(2, float64(z))
	return t
}

func (t *TileCoordinateBound) TileBounds() *TileCoordinateBound {
	tile3857Miny := math.Pow(2, float64(t.zoom)) - 1 - float64(t.MaxY)
	tile3857Maxy := math.Pow(2, float64(t.zoom)) - 1 - float64(t.MinY)

	minx, miny := t.PixelsToMeters(float64(t.MinX*t.tileSize), tile3857Miny*float64(t.tileSize), t.zoom)
	maxx, maxy := t.PixelsToMeters(float64((t.MaxX+1)*t.tileSize), (tile3857Maxy+1)*float64(t.tileSize), t.zoom)
	t.rMinX = minx
	t.rMinY = miny
	t.rMaxX = maxx
	t.rMaxY = maxy
	return t
}

func (t *TileCoordinateBound) PixelsToMeters(px, py float64, zoom int) (mx, my float64) {
	if t.resolution == 0 {
		t.Resolution(zoom)
	}

	mx = math.Abs(float64(px)*t.resolution - t.originShift)
	my = math.Abs(py*t.resolution - t.originShift)
	return
}

func (t *TileCoordinateBound) GetResolutionBounds() struct{ MinX, MinY, MaxX, MaxY float64 } {
	if t.rMinY == 0 && t.rMinX == 0 && t.rMaxY == 0 && t.rMaxX == 0 {
		t.TileBounds()
	}
	return struct {
		MinX, MinY, MaxX, MaxY float64
	}{
		MinX: t.rMinX,
		MinY: t.rMinY,
		MaxX: t.rMaxX,
		MaxY: t.rMaxY,
	}
}

func (t *TileCoordinateBound) GetResolution() float64 {
	if t.resolution == 0 {
		t.Resolution(t.zoom)
	}
	return t.resolution
}

func (t *TileCoordinateBound) GetGeoTransform() [6]float64 {
	if t.resolution == 0 {
		t.Resolution(t.zoom)
	}
	if t.rMinY == 0 && t.rMinX == 0 && t.rMaxY == 0 && t.rMaxX == 0 {
		t.TileBounds()
	}

	return [6]float64([]float64{t.rMinX, t.resolution, 0.000000000000000, t.rMaxY, 0.000000000000000, -t.resolution})
}
