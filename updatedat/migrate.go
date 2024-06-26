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
				Dialect: sqlbuilder.Dialect_mysql,
				Scene:   sqlbuilder.SCENE_DDL_CREATE,
				DDL:     fmt.Sprintf("`%s` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',", Field_UpdatedAt.Name),
			},
			{
				Dialect: sqlbuilder.Dialect_mysql,
				Scene:   sqlbuilder.SCENE_DDL_APPEND,
				DDL:     fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间' %s;", table, Field_CreatedAt.Name, mysqlAfter.String()),
			},
			{
				Dialect: sqlbuilder.Dialect_mysql,
				Scene:   sqlbuilder.SCENE_DDL_MODIFY,
				DDL:     fmt.Sprintf("ALTER TABLE `%s` MODIFY `%s` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',", table, Field_UpdatedAt.Name),
			},
			{
				Dialect: sqlbuilder.Dialect_mysql,
				Scene:   sqlbuilder.SCENE_DDL_DELETE,
				DDL:     fmt.Sprintf("ALTER TABLE `%s` DROP `%s` ;", table, Field_UpdatedAt.Name),
			},
		}
	}
}

var Field_CreatedAt = sqlbuilder.Field{
	Name:  "created_at",
	Title: "创建时间",
	Value: func(in any) (any, error) { return in, nil },
}

var Field_UpdatedAt = sqlbuilder.Field{
	Name:  "updated_at",
	Title: "更新时间",
	Value: func(in any) (any, error) { return in, nil },
}

func Migrate(table string, driver sqlbuilder.Driver, scene sqlbuilder.Scene) []string {
	all := make(sqlbuilder.Migrates, 0)
	all = append(all, Field_CreatedAt.Migrate(table)...)
	all = append(all, Field_UpdatedAt.Migrate(table)...)
	sub := all.GetByScene(driver, scene)
	return sub.DDLs()
}
