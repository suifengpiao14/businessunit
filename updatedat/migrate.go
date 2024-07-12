package updatedat

import (
	"fmt"

	"github.com/suifengpiao14/sqlbuilder"
)

func init() {
	Field_UpdatedAt.Migrate = func(table string, options ...sqlbuilder.MigrateOptionI) sqlbuilder.Migrates {
		mysqlAfter := sqlbuilder.GetMigrateOpion(sqlbuilder.MigrateOptionMysqlAfter(""), options...)
		return sqlbuilder.Migrates{
			{
				Dialect: sqlbuilder.Driver_mysql,
				Scene:   sqlbuilder.SCENE_DDL_CREATE,
				DDL:     fmt.Sprintf("`%s` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',", sqlbuilder.FieldName2DBColumnName),
			},
			{
				Dialect: sqlbuilder.Driver_mysql,
				Scene:   sqlbuilder.SCENE_DDL_APPEND,
				DDL:     fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间' %s;", table, sqlbuilder.FieldName2DBColumnName, mysqlAfter.String()),
			},
			{
				Dialect: sqlbuilder.Driver_mysql,
				Scene:   sqlbuilder.SCENE_DDL_MODIFY,
				DDL:     fmt.Sprintf("ALTER TABLE `%s` MODIFY `%s` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',", table, sqlbuilder.FieldName2DBColumnName),
			},
			{
				Dialect: sqlbuilder.Driver_mysql,
				Scene:   sqlbuilder.SCENE_DDL_DELETE,
				DDL:     fmt.Sprintf("ALTER TABLE `%s` DROP `%s` ;", table, sqlbuilder.FieldName2DBColumnName),
			},
		}
	}
}

var Field_CreatedAt = sqlbuilder.NewField(func(in any) (any, error) { return in, nil }).SetName("created_at")

var Field_UpdatedAt = sqlbuilder.NewField(func(in any) (any, error) { return in, nil }).SetName("updated_at")

func Migrate(table string, driver sqlbuilder.Driver, scene sqlbuilder.Scene) []string {
	all := make(sqlbuilder.Migrates, 0)
	all = append(all, Field_CreatedAt.Migrate(table)...)
	all = append(all, Field_UpdatedAt.Migrate(table)...)
	sub := all.GetByScene(driver, scene)
	return sub.DDLs()
}
