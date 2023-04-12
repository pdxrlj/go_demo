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
	"image"
	"image/png"
	"math"
	"os"
	"path/filepath"

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
	vrt := dataset.mercatorVrt()
	for i := 7; i < 11; i++ {
		tileRange := NewTileRange().TileRange(i, bounds)
		tileRange.iter(func(tileId *TileId) {
			join := filepath.Join(fmt.Sprintf("tmp/%d/%d/%d.png", tileId.z, tileId.x, tileId.y))

			basepath := filepath.Dir(join)
			if _, err := os.Stat(basepath); os.IsNotExist(err) {
				err := os.MkdirAll(basepath, 0755)
				if err != nil {
					panic(err)
				}
			}
			buf := vrt.ReadTile(tileId, 256)

			img := &image.Gray{
				Pix:    buf,
				Stride: 256,
				Rect:   image.Rect(0, 0, 256, 256),
			}

			f, _ := os.Create(join)
			err = png.Encode(f, img)
			if err != nil {
				return
			}
			return
		})
	}
	vrt.vrtDataset.Close()
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
	invDeterminant := 1.0 / (a.a*a.e - a.b*a.d)

	ac := &Affine{}
	ac.a = a.e * invDeterminant
	ac.b = -a.b * invDeterminant
	ac.d = -a.d * invDeterminant
	ac.e = a.a * invDeterminant

	ac.c = -a.c*ac.a - a.f*ac.b
	ac.f = -a.c*ac.d - a.f*ac.e
	return ac
}

func (a *Affine) multiply(x, y float64) (float64, float64) {
	xc := x*a.a + y*a.b + a.c
	yc := x*a.d + y*a.e + a.f
	return xc, yc
}

func (a *Affine) scale(x, y float64) *Affine {
	return &Affine{
		a: a.a * x,
		b: a.b,
		c: a.c,
		d: a.d,
		e: a.e * y,
		f: a.f,
	}
}

func (a *Affine) resolution() (float64, float64) {
	return math.Abs(a.a), math.Abs(a.e)
}

type Windows struct {
	xOffset float64
	yOffset float64
	width   float64
	height  float64
}

func NewWindows() *Windows {
	return &Windows{}
}

func (w *Windows) fromBounds(transform *Affine, bounds *Bounds) *Windows {
	transformInvert := transform.invert()

	xs := [4]float64{0.0, 0.0, 0.0, 0.0}
	ys := [4]float64{0.0, 0.0, 0.0, 0.0}
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

	return &Windows{
		xOffset: xmin,
		yOffset: ymin,
		width:   xmax - xmin,
		height:  ymax - ymin,
	}
}

func (w *Windows) transform(transform *Affine) *Affine {
	x, y := transform.multiply(w.xOffset, w.yOffset)
	return &Affine{
		a: transform.a,
		b: transform.b,
		c: x,
		d: transform.d,
		e: transform.e,
		f: y,
	}
}

// -----

type Dataset struct {
	dataset       *godal.Dataset
	vrtDataset    gdal.Dataset
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

func (d *Dataset) ReadTile(tileId *TileId, tileSize int) []byte {

	vrtWidthF := float64(d.vrtDataset.RasterXSize())
	vrtHeightF := float64(d.vrtDataset.RasterYSize())

	geoTransform := d.vrtDataset.GeoTransform()

	affine := NewAffine().fromGdal(geoTransform)

	vrtBounds := d.Bounds()

	tileBounds := tileId.MercatorBounds()

	windows := NewWindows().fromBounds(affine, tileBounds)

	tileTransform := windows.transform(affine).scale(windows.width/256, windows.height/256)

	xres, yres := tileTransform.resolution()

	left := math.Max(0, math.Round((vrtBounds.xmin-tileBounds.xmin)/xres))
	right := math.Max(0, math.Round((tileBounds.xmax-vrtBounds.xmax)/xres))
	bottom := math.Max(0, math.Round((vrtBounds.ymin-tileBounds.ymin)/yres))
	top := math.Max(0, math.Round((tileBounds.ymax-vrtBounds.ymax)/yres))

	widthSize := int(math.Round(float64(tileSize) - left - right))
	heightSize := int(math.Round(float64(tileSize) - top - bottom))

	xOffset := math.Round(math.Min(vrtWidthF, math.Max(0, windows.xOffset)))
	yOffset := math.Round(math.Min(vrtHeightF, math.Max(0, windows.yOffset)))
	xStop := math.Min(vrtWidthF, math.Max(0, windows.xOffset+windows.width))
	yStop := math.Min(vrtHeightF, math.Max(0, windows.yOffset+windows.height))

	readWidth := int(math.Floor(float64(xStop-xOffset) + 0.5))
	readHeight := int(math.Floor(float64(yStop-yOffset) + 0.5))
	if readWidth <= 0 || readHeight <= 0 {
		return nil
	}

	//span := widthsize * heightsize * d.vrtDataset.RasterCount()
	//tmp := make([]byte, span, span)
	//err := d.vrtDataset.IO(gdal.Read, xOffset, yOffset, readWidth, readHeight, tmp, widthsize, heightsize, d.vrtDataset.RasterCount(), []int{1, 2, 3}, 0, 0, 0)
	//if err != nil {
	//	panic(err)
	//}

	//span := widthSize * heightSize
	tmp := make([]byte, 256*256, 256*256)

	for i := range tmp {
		tmp[i] = 255
	}

	err := d.vrtDataset.RasterBand(1).IO(gdal.Read, int(xOffset), int(yOffset), readWidth, readHeight, tmp, widthSize, heightSize, 0, 0)
	if err != nil {
		panic(err)
	}

	if left > 0 || top > 0 || widthSize < tileSize || heightSize < tileSize {
		buff := shift(tileId, tmp, [2]int{widthSize, heightSize}, [2]int{tileSize, tileSize}, [2]int{int(left), int(top)}, 255)
		tmp = buff
	}

	return tmp
	//create, err := godal.Create(godal.GTiff, join, d.vrtDataset.RasterCount(), godal.Byte, 256, 256)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("tmplen:%d widthSpan:%d heightsize:%d\n", len(tmp), widthsize*heightsize, heightsize)
	//err = create.Write(0, 0, tmp, span*3, span*3)
	//if err != nil {
	//	panic(err)
	//}
	//defer create.Close()
	//utils.WriteFinal(heightsize, widthsize, tmp, join)
	//utils.Write(3, widthsize, heightsize, tmp, join)

	//utils.Write2(3, xOffset, yOffset, readWidth, readHeight, widthsize, heightsize, join, d.vrtDataset)

}

func shift(tileId *TileId, buffer []byte, size, targetSize, offset [2]int, fill byte) []byte {
	var srcIndex, destIndex int

	for row := size[1] - 1; row >= 0; row-- {
		for col := size[0] - 1; col >= 0; col-- {
			srcIndex = row*size[0] + col
			destIndex = (row+offset[1])*targetSize[0] + (col + offset[0])
			buffer[destIndex] = buffer[srcIndex]

			if destIndex != srcIndex {
				buffer[srcIndex] = fill
			}
		}
	}
	return buffer
}

func (d *Dataset) Bounds() *Bounds {
	w := d.vrtDataset.RasterXSize()
	h := d.vrtDataset.RasterYSize()
	geoTransform := d.vrtDataset.GeoTransform()
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
	vrt := d.warpedVrt()
	return &Dataset{
		vrtDataset:    vrt,
		originDataSet: d.originDataSet,
	}
}

func (d *Dataset) warpedVrt() gdal.Dataset {
	dataset, err := gdal.Open(d.originDataSet, gdal.ReadOnly)
	if err != nil {
		panic(err)
	}
	wkt, err := d.dataset.SpatialRef().WKT()
	if err != nil {
		panic(err)
	}
	reference, _ := godal.NewSpatialRefFromEPSG(3857)
	dst, err := reference.WKT()
	if err != nil {
		panic(err)
	}

	vrt, err := dataset.AutoCreateWarpedVRT(wkt, dst, gdal.GRA_NearestNeighbour)
	if err != nil {
		panic(err)
	}
	//vrtWarp, err := gdal.Warp("demo_wrap.vrt",
	//	nil,
	//	[]gdal.Dataset{dataset}, []string{
	//		"-r", "near",
	//		"-t_srs", "EPSG:3857",
	//		"-of", "VRT",
	//		"-wm", "512",
	//		"-multi", "2",
	//	})
	if err != nil {
		panic(err)
	}

	return vrt
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
