package polygon

import (
	"github.com/suifengpiao14/sqlbuilder"
)

func (boundingBox BoundingBox) Fields() (boundingBoxFields sqlbuilder.Fields) {
	var (
		lngMaxName = "lngMax"
		lngMinName = "lngMin"
		latMaxName = "latMax"
		latMinName = "latMin"
	)
	boundingBoxFields = sqlbuilder.Fields{
		sqlbuilder.NewField(func(in any) (any, error) { return boundingBox.LngMax, nil }).SetName(lngMaxName).SetTitle("最大经度").SetTag(lngMaxName),
		sqlbuilder.NewField(func(in any) (any, error) { return boundingBox.LngMin, nil }).SetName(lngMinName).SetTitle("最小经度").SetTag(lngMinName),
		sqlbuilder.NewField(func(in any) (any, error) { return boundingBox.LngMax, nil }).SetName(latMaxName).SetTitle("最大纬度").SetTag(latMaxName),
		sqlbuilder.NewField(func(in any) (any, error) { return boundingBox.LngMin, nil }).SetName(latMinName).SetTitle("最小纬度").SetTag(latMinName),
	}
	boundingBoxFields.SceneSelect(func(f *sqlbuilder.Field, fs ...*sqlbuilder.Field) {
		fields := sqlbuilder.Fields(fs)
		if f.HastTag(latMaxName) {
			latMinField, _ := fields.GetByTag(latMinName)
			f.WhereFns.Append(sqlbuilder.ValueFn{
				Layer: sqlbuilder.Value_Layer_SetValue,
				Fn: func(data any) (any, error) {
					return sqlbuilder.Between{latMinField.DBName(), data, f.DBName()}, nil
				},
			})
		} else if f.HastTag(latMinName) {
			f.WhereFns.Append(sqlbuilder.ValueFnShield)
		} else if f.HastTag(lngMaxName) {
			LngMinField, _ := fields.GetByTag(lngMinName)
			f.WhereFns.Append(sqlbuilder.ValueFn{
				Layer: sqlbuilder.Value_Layer_SetValue,
				Fn: func(data any) (any, error) {
					return sqlbuilder.Between{LngMinField.DBName(), data, f.DBName()}, nil
				},
			})
		} else if f.HastTag(lngMinName) {
			f.WhereFns.Append(sqlbuilder.ValueFnShield)

		}
	})
	return boundingBoxFields
}
