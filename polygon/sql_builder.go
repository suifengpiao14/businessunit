package polygon

import (
	"github.com/suifengpiao14/sqlbuilder"
)

type BoundingBoxField struct {
	LngMax *sqlbuilder.Field
	LngMin *sqlbuilder.Field
	LatMax *sqlbuilder.Field
	LatMin *sqlbuilder.Field
}

func (boundingBoxField BoundingBoxField) Fields() sqlbuilder.Fields {
	return sqlbuilder.Fields{
		boundingBoxField.LatMax,
		boundingBoxField.LatMin,
		boundingBoxField.LngMax,
		boundingBoxField.LngMin,
	}
}

type PolygonI interface {
	Points() (points Points, err error)
	GetBoundingBoxField() (boundingBoxField BoundingBoxField)
}

func SetBoundingBoxFieldValue(polygonI PolygonI) (er error) {
	points, err := polygonI.Points()
	if err != nil {
		return err
	}
	if len(points) == 0 {
		return
	}
	boundingBox, err := points.GetBoundingBox()
	if err != nil {
		return err
	}
	boundingBoxField := polygonI.GetBoundingBoxField()
	boundingBoxField.LngMax.ValueFns.InsertAsFirst(func(in any) (any, error) {
		return boundingBox.LngMax, nil
	})
	boundingBoxField.LngMin.ValueFns.InsertAsFirst(func(in any) (any, error) {
		return boundingBox.LngMin, nil
	})
	boundingBoxField.LatMax.ValueFns.InsertAsFirst(func(in any) (any, error) {
		return boundingBox.LatMax, nil
	})
	boundingBoxField.LatMin.ValueFns.InsertAsFirst(func(in any) (any, error) {
		return boundingBox.LatMin, nil
	})
	return nil
}

func Insert(polygonI PolygonI) (err error) {
	return SetBoundingBoxFieldValue(polygonI)
}

func Update(polygonI PolygonI) (err error) {
	return SetBoundingBoxFieldValue(polygonI)
}

func Select(boundingBox BoundingBoxField) {
	boundingBox.LatMax.WhereFns.AppendIfNotFirst(func(data any) (any, error) {
		if sqlbuilder.IsNil(data) {
			return nil, nil
		}
		return sqlbuilder.Between{boundingBox.LatMin.DBName(), boundingBox.LatMax.DBName()}, nil
	})
	boundingBox.LatMax.WhereFns.Append(sqlbuilder.ValueFnShield)

	boundingBox.LngMax.WhereFns.AppendIfNotFirst(func(data any) (any, error) {
		if sqlbuilder.IsNil(data) {
			return nil, nil
		}
		return sqlbuilder.Between{boundingBox.LngMin.DBName(), boundingBox.LngMax.DBName()}, nil
	})
	boundingBox.LngMax.WhereFns.Append(sqlbuilder.ValueFnShield)
}
