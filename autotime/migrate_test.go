package autotime_test

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/businessunit/autotime"
	"github.com/suifengpiao14/sqlbuilder"
)

func TestMigrate(t *testing.T) {
	table := "t_user"
	autotime.Field_CreatedAt.Name = "Fauto_create_time"
	autotime.Field_UpdatedAt.Name = "Fauto_update_time"
	ddls := autotime.Migrate(table, sqlbuilder.Dialect_mysql, sqlbuilder.SCENE_APPEND)
	fmt.Println(ddls)

}
