// 基础概念
// resolution每像素的地图单位 【可以理解为每像素代表多少米】
// resolution=全幅米 / 全幅像素
// web墨卡托切片瓦片的起始范围的原点是坐标为（-20037508.342789244,20037508.34278924）
// 投影后地图最大范围-20037508.342789244, 20037508.342789244, 20037508.342789244, -20037508.342789244,
// 赤道半径为6378137米

// transform
// geos[0]top left x 左上角x坐标
// geos[1] w-e pixel resolution 东西方向像素分辨率
// geos[2] rotation，0 if image is "north up" 旋转角度，正北向上时为0
// geos[3] top left y 左上角y坐标
// geos[4] rotation，@ if image is "north up" 旋转角度，正北向上时为0
// geos[5] n-s pixel resolution 南北向像素分辨率

package main

import (
	"fmt"
	"math"

	"github.com/airbusgeo/godal"
	"github.com/lukeroth/gdal"
)

func main() {
	dataset := NewDataset("/Users/ruanlianjun/Desktop/demo.tif")
	epsg, err := godal.NewSpatialRefFromEPSG(3857)
	if err != nil {
		panic(err)
	}

	bounds := dataset.TransformBounds(epsg)
	tileRange := NewTileRange().TileRange(10, bounds)

	vrt := dataset.mercatorVrt()
	// vrt_width:11314, vrt_height:11365
	fmt.Printf("width:%d height:%d\n", vrt.dataset.Structure().SizeX, vrt.dataset.Structure().SizeY)
	tileRange.iter(func(tileId *TileId) {
		//tmp := make([]byte, 256*256*3)
		//vrt.ReadTile(tileId, 256, tmp)
		return
	})

}

type TileRange struct {
	zoom int
	xmin int
	ymin int
	xmax int
	ymax int
}

func NewTileRange() *TileRange {
	return &TileRange{}
}

func (t *TileRange) TileRange(zoom int, bounds *Bounds) *TileRange {
	zInt := 1 << zoom
	z := float64(zInt)
	origin := -ORIGIN
	eps := 1e-11

	xmin := (bounds.xmin - origin) / CE * z
	xmin = math.Max(xmin, 0.0)
	xmin = math.Min(xmin, z-1.0)
	xmin = math.Floor(xmin)
	xmin = float64(uint32(xmin))

	ymin := (1.0 - (bounds.ymax-origin)/CE) * z
	ymin = math.Max(ymin, 0.0)
	ymin = math.Min(ymin, z-1.0)
	ymin = math.Floor(ymin)
	ymin = float64(uint32(ymin))

	xmax := ((bounds.xmax-origin)/CE - eps) * z
	xmax = math.Max(xmax, 0.0)
	xmax = math.Min(xmax, z-1.0)
	xmax = math.Floor(xmax)
	xmax = float64(uint32(xmax))

	ymax := (1.0 - (bounds.ymin-origin)/CE + eps) * z
	ymax = math.Max(ymax, 0.0)
	ymax = math.Min(ymax, z-1.0)
	ymax = math.Floor(ymax)
	ymax = float64(uint32(ymax))

	tileRange := TileRange{
		zoom: zoom,
		xmin: int(uint32(xmin)),
		ymin: int(uint32(ymin)),
		xmax: int(uint32(xmax)),
		ymax: int(uint32(ymax)),
	}

	return &tileRange
}

func (t *TileRange) iter(fn func(tileId *TileId)) {
	for i := t.xmin; i <= t.xmax; i++ {
		for j := t.ymin; j <= t.ymax; j++ {
			tile := NewTileID(t.zoom, i, j)
			fn(tile)
		}
	}
}

//--------

type Affine struct {
	a float64
	b float64
	c float64
	d float64
	e float64
	f float64
}

func NewAffine() *Affine {
	return &Affine{}
}

func (a *Affine) fromGdal(transform [6]float64) *Affine {
	a.a = transform[1]
	a.b = transform[2]
	a.c = transform[0]
	a.d = transform[4]
	a.e = transform[5]
	a.f = transform[3]
	return a
}

func (a *Affine) invert() *Affine {
	inv_determinant := 1.0 / (a.a*a.e - a.b*a.d)

	ac := &Affine{}
	ac.a = a.e * inv_determinant
	ac.b = -a.b * inv_determinant
	ac.d = -a.d * inv_determinant
	ac.e = a.a * inv_determinant
	ac.c = -a.c*a.a - a.f*a.b
	ac.f = -a.c*a.d - a.f*a.e
	return ac
}

func (a *Affine) multiply(x, y float64) (float64, float64) {
	// x * self.a + y * self.b + self.c,
	// x * self.d + y * self.e + self.f,
	xc := x*a.a + y*a.b + a.c
	yc := x*a.d + y*a.e + a.f
	return xc, yc
}

type Windows struct {
	x_offset float64
	y_offset float64
	width    float64
	height   float64
}

func NewWindows() *Windows {
	return &Windows{}
}

func (w *Windows) fromBounds(transform *Affine, bounds *Bounds) *Windows {
	transformInvert := transform.invert()
	xs := [4]float64{0, 0, 0, 0}
	ys := [4]float64{0, 0, 0, 0}
	xs[0], ys[0] = transformInvert.multiply(bounds.xmin, bounds.ymin)
	xs[1], ys[1] = transformInvert.multiply(bounds.xmin, bounds.ymax)
	xs[2], ys[2] = transformInvert.multiply(bounds.xmax, bounds.ymin)
	xs[3], ys[3] = transformInvert.multiply(bounds.xmax, bounds.ymax)
	xmin := xs[0]
	xmax := xs[0]
	ymin := ys[0]
	ymax := ys[0]
	for i := 1; i < 4; i++ {
		if xs[i] < xmin {
			xmin = xs[i]
		}
		if ys[i] < ymin {
			ymin = ys[i]
		}
		if xs[i] > xmax {
			xmax = xs[i]
		}
		if ys[i] > ymax {
			ymax = ys[i]
		}
	}

	return w
}

// -----

type Dataset struct {
	dataset       *godal.Dataset
	originDataSet string
}

func NewDataset(path string) *Dataset {
	err := godal.RegisterRaster("GTiff")
	if err != nil {
		panic(err)
	}
	dataset, err := godal.Open(path)
	if err != nil {
		panic(err)
	}
	return &Dataset{
		originDataSet: path,
		dataset:       dataset,
	}

}

func (d *Dataset) ReadTile(tileId *TileId, tileSize int, buffer []byte) {
	//tileSizeF := float64(tileSize)

	// tile_id:TileID { zoom: 10, x: 833, y: 405 } vrt_width:11314, vrt_height:11365
	// windows:Window { x_offset: -2936.607967421529, y_offset: 554.7919124158216, width: 3231.292925737798, height: 3231.292925737798 }

	vrtWidthF := d.dataset.Structure().SizeX
	vrtHeightF := d.dataset.Structure().SizeY
	fmt.Printf("tileId:%+v vrt_width_f:%+v vrt_width_f:%+v\n", tileId, vrtWidthF, vrtHeightF)
	return
	geoTransform, err := d.dataset.GeoTransform()
	if err != nil {
		panic(err)
	}
	affine := NewAffine().fromGdal(geoTransform)

	vrtBounds := d.Bounds()
	fmt.Printf("vrt_bounds:%+v\n", vrtBounds)

	tileBounds := tileId.MercatorBounds()
	//{ xmin: 12598145.158510394, ymin: 4056598.475505117, xmax: 12735174.509188766, ymax: 4194245.511960808 }
	fmt.Printf("tile_bounds:%+v\n", tileBounds)

	windows := NewWindows().fromBounds(affine, tileBounds)
	fmt.Printf("windows:%+v\n", windows)

}

func (d *Dataset) Bounds() *Bounds {
	w := d.dataset.Structure().SizeX
	h := d.dataset.Structure().SizeY
	geoTransform, err := d.dataset.GeoTransform()
	if err != nil {
		panic(err)
	}
	affine := NewAffine().fromGdal(geoTransform)

	return &Bounds{
		xmin: affine.c,
		ymin: affine.f + affine.e*float64(h),
		xmax: affine.c + affine.a*float64(w),
		ymax: affine.f,
	}
}

func (d *Dataset) mercatorBounds() *Bounds {
	epsg, err := godal.NewSpatialRefFromEPSG(3857)
	if err != nil {
		panic(err)
	}
	bounds := d.TransformBounds(epsg)
	return bounds
}

func (d *Dataset) TransformBounds(crs *godal.SpatialRef) *Bounds {
	bounds, err := d.dataset.Bounds()
	if err != nil {
		panic(err)
	}
	src := d.dataset.SpatialRef()
	transform, err := godal.NewTransform(src, crs)
	x := []float64{bounds[0], bounds[2]}
	y := []float64{bounds[1], bounds[3]}
	if err = transform.TransformEx(x, y, nil, nil); err != nil {
		panic(err)
	}
	transform.Close()

	return &Bounds{
		xmin: x[0],
		ymin: y[0],
		xmax: x[1],
		ymax: y[1],
	}
}

func (d *Dataset) mercatorVrt() *Dataset {
	epsg, _ := godal.NewSpatialRefFromEPSG(3857)
	vrt := d.warpedVrt(epsg)
	return &Dataset{
		dataset:       vrt,
		originDataSet: d.originDataSet,
	}
}

func (d *Dataset) warpedVrt(dst *godal.SpatialRef) *godal.Dataset {
	epsg, err := godal.NewSpatialRefFromEPSG(3857)
	if err != nil {
		panic(err)
	}
	dstSpa, _ := epsg.WKT()
	dataset, err := godal.BuildVRT("demo.vrt", []string{d.originDataSet},
		[]string{
			"-r", "nearest",
			"-a_srs", dstSpa,
		},
		godal.ConfigOption(
			"GDAL_NUM_THREADS=1",
			"NUM_THREADS=1",
			"eResampleAlg=nearest"),
		godal.Resampling(godal.Nearest))
	if err != nil {
		panic(err)
	}
	open, err := gdal.Open(d.originDataSet, gdal.ReadOnly)
	if err != nil {
		panic(err)
	}

	spatialRef := d.dataset.SpatialRef()
	wkt, err := spatialRef.WKT()
	if err != nil {
		panic(err)
	}
	vrt, err := gdal.BuildVRT(
		"",
		[]gdal.Dataset{open},
		[]string{d.originDataSet},
		[]string{
			"srcSRS", wkt,
			"VRTSRS", dstSpa,
			"resampling", "NearestNeighbor",
		})
	if err != nil {
		panic(err)
	}

	fmt.Printf("vrt:%+v\n", vrt.RasterXSize())
	fmt.Printf("vrt:%+v\n", vrt.RasterYSize())

	return dataset
}

// -----------

type Bounds struct {
	xmin float64
	ymin float64
	xmax float64
	ymax float64
}

type TileId struct {
	z, x, y int
}

func NewTileID(z, x, y int) *TileId {
	return &TileId{
		z: z,
		x: x,
		y: y,
	}
}

const RE = 6378137.0
const ORIGIN = RE * math.Pi
const CE = 2.0 * ORIGIN

func (tile *TileId) MercatorBounds() *Bounds {
	z := 1 << tile.z
	x := float64(tile.x)
	y := float64(tile.y)

	tileSize := CE / float64(z)
	xmin := x*tileSize - CE/2.0
	ymax := CE/2.0 - y*tileSize
	return &Bounds{
		xmin: xmin,
		ymin: ymax - tileSize,
		xmax: xmin + tileSize,
		ymax: ymax,
	}
}

func (tile TileId) GeoBounds() *Bounds {
	rad2deg := 180.0 / math.Pi
	zi := 1 << tile.z
	z := float64(zi)
	x := float64(tile.x)
	y := float64(tile.y)
	return &Bounds{
		xmin: x/z*360.0 - 180.0,
		ymin: math.Sinh(math.Atan(math.Sinh(math.Pi*(1.0-2.0*((y+1.0)/z))))) * rad2deg,
		xmax: (x+1.0)/z*360.0 - 180.0,
		ymax: math.Sinh(math.Atan(math.Sinh(math.Pi*(1.0-2.0*y/z)))) * rad2deg,
	}
}

// Resolution 根据层级计算分辨率
func Resolution(z int) float64 {
	// 列数
	cols := math.Pow(2.0, float64(z))
	// 全幅像素
	var pixels float64 = float64(cols) * 256.0
	// 全幅米
	meters := 20037508.342789244 + 20037508.342789244 //40075016.68557849;
	// 当前层级每像素的地图单位
	resolution := meters / pixels

	return resolution
}

const (
	R           = 6378137
	MaxLatitude = 85.0511287798
)

func LonLatToMercator(lon, lat float64) (float64, float64) {
	d := math.Pi / 180
	max := MaxLatitude
	lat = math.Max(math.Min(max, lat), -max)
	sin := math.Sin(lat * d)
	x := R * lon * d
	y := R * math.Log((1+sin)/(1-sin)) / 2
	return x, y
}

func MercatorToLongLat(x, y float64) (float64, float64) {
	d := 180 / math.Pi
	lon := x * d / R
	lat := (2*math.Atan(math.Exp(y/R)) - (math.Pi / 2)) * d
	return lon, lat
}

// ComputeStartXY 计算瓦片起始坐标
// @param startX 起始x坐标 米
// @param startY 起始y坐标 米
func ComputeStartXY(startX, startY float64, z int) (float64, float64) {
	//当前位置x方向距离原点多少米
	x := startX - (-20037508.342789244)
	//当前位置x方向距离原点多少像素
	xPixel := x / Resolution(z)
	//瓦片横坐标
	xTile := math.Floor(xPixel / 256)

	//当前位置y方向距离原点多少米
	y := 20037508.342789244 - startY
	// 当前位置y方向距离原点多少像素
	yPixel := y / Resolution(z)
	// 瓦片纵坐标
	yTile := math.Floor(yPixel / 256)
	return xTile, yTile
}

func ComputeResolution(z int, originWidth, originHeight, transform1, transform5 float64) (float64, float64) {
	resolution := Resolution(z)
	x := math.Round(originWidth * transform1 / resolution)
	y := math.Round(originHeight * transform5 / resolution)

	return x, y
}

func MinMaxZoom(xMinLon, xMaxLong, xSize, ySize float64) (int, int) {
	minzoom := int(math.Log2(360 / (xMaxLong - xMinLon)))

	maxsize := int(math.Max(xSize/256, ySize/256))
	levels := int(math.Ceil(math.Log2(float64(maxsize))))
	maxzoom := minzoom + levels - 1
	return minzoom, maxzoom
}
