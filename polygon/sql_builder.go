package polygon

import (
	"github.com/suifengpiao14/sqlbuilder"
)

type BoundingBoxField struct {
	LngMax sqlbuilder.Field
	LngMin sqlbuilder.Field
	LatMax sqlbuilder.Field
	LatMin sqlbuilder.Field
}

type PolygonI interface {
	Points() (points Points, err error)
	GetBoundingBoxField() (boundingBoxField BoundingBoxField)
}

func mergePolygonData(polygon PolygonI) (polygonData map[string]any, err error) {
	points, err := polygon.Points()
	if err != nil {
		return nil, err
	}

	polygonData = make(map[string]interface{})
	if len(points) > 0 {
		boundingBox, err := points.GetBoundingBox()
		if err != nil {
			return nil, err
		}
		boundingBoxField := polygon.GetBoundingBoxField()
		lngMax, err := boundingBoxField.LngMax.GetValue(boundingBox.LngMax)
		if err != nil {
			return nil, err
		}
		lngMin, err := boundingBoxField.LngMin.GetValue(boundingBox.LngMin)
		if err != nil {
			return nil, err
		}
		latMax, err := boundingBoxField.LatMax.GetValue(boundingBox.LatMax)
		if err != nil {
			return nil, err
		}
		latMin, err := boundingBoxField.LatMin.GetValue(boundingBox.LatMin)
		if err != nil {
			return nil, err
		}
		polygonData[sqlbuilder.FieldName2DBColumnName(boundingBoxField.LngMax.Name)] = lngMax
		polygonData[sqlbuilder.FieldName2DBColumnName(boundingBoxField.LngMin.Name)] = lngMin
		polygonData[sqlbuilder.FieldName2DBColumnName(boundingBoxField.LatMax.Name)] = latMax
		polygonData[sqlbuilder.FieldName2DBColumnName(boundingBoxField.LatMin.Name)] = latMin
	}

	return polygonData, nil
}

type _Polygon struct {
	PolygonI
}

func (t _Polygon) Data() (data interface{}, err error) {
	return mergePolygonData(t)
}

func NewPolygon(polygonI PolygonI) _Polygon {
	return _Polygon{
		PolygonI: polygonI,
	}
}

func Insert(param PolygonI) sqlbuilder.InsertParam {
	return sqlbuilder.NewInsertBuilder(nil).AppendData(NewPolygon(param))
}

func Update(param PolygonI) sqlbuilder.UpdateParam {
	return sqlbuilder.NewUpdateBuilder(nil).AppendData(NewPolygon(param))
}

/**以下仅仅为了完备,方便调用方使用,减少调用方心智负担, 统一格式后,也方便调用方批量处理中间件**/

func First(param PolygonI) sqlbuilder.FirstParam {
	return sqlbuilder.NewFirstBuilder(nil)
}

func List(param PolygonI) sqlbuilder.ListParam {
	return sqlbuilder.NewListBuilder(nil)
}

func Total(param PolygonI) sqlbuilder.TotalParam {
	return sqlbuilder.NewTotalBuilder(nil)
}
