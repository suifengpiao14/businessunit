package polygon

import (
	"math"

	geo "github.com/kellydunn/golang-geo"
)

type BoundingBox struct {
	LatMin float64
	LatMax float64
	LngMin float64
	LngMax float64
}

func (box BoundingBox) ContainsPoint(targetPoint *geo.Point) (ok bool) {
	ok = (box.LngMin <= targetPoint.Lng() && targetPoint.Lng() <= box.LngMax) && (box.LatMin <= targetPoint.Lat() && targetPoint.Lat() <= box.LatMax)
	return
}

type Point struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

type Points []Point

// geoPolygon 数据表中的点转换为geo.Polygon,这个函数和数据存储格式相关，抽象的话需要能重写
func (ps Points) geoPolygon() (polygon *geo.Polygon, err error) {
	geoPointSet := make([]*geo.Point, 0)
	for _, point := range ps {
		geoPint := geo.NewPoint(point.Lat, point.Lng)
		geoPointSet = append(geoPointSet, geoPint)
	}
	polygon = geo.NewPolygon(geoPointSet)
	return polygon, nil
}

// GetBoundingBox 获取多边形 Bounding Box
func (ps Points) GetBoundingBox() (boundingBox *BoundingBox, err error) {
	minLat := math.MaxFloat64
	maxLat := -math.MaxFloat64
	minLng := math.MaxFloat64
	maxLng := -math.MaxFloat64
	for _, point := range ps {
		if point.Lat < minLat {
			minLat = point.Lat
		}
		if point.Lat > maxLat {
			maxLat = point.Lat
		}
		if point.Lng < minLng {
			minLng = point.Lng
		}
		if point.Lng > maxLng {
			maxLng = point.Lng
		}
	}
	boundingBox = &BoundingBox{
		LatMin: minLat,
		LatMax: maxLat,
		LngMin: minLng,
		LngMax: maxLng,
	}
	return boundingBox, nil
}

// ContainsPoint 判断点是否在多边形内
func (ps Points) ContainsPoint(targetPoint *geo.Point) (ok bool, err error) {
	boundingBox, err := ps.GetBoundingBox()
	if err != nil {
		return false, err
	}
	ok = boundingBox.ContainsPoint(targetPoint)
	if !ok {
		return false, nil
	}
	geoPolygon, err := ps.geoPolygon()
	if err != nil {
		return false, err
	}
	ok = geoPolygon.Contains(targetPoint)
	return ok, nil
}
