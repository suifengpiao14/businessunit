package updatedat_test

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/businessunit/updatedat"
	"github.com/suifengpiao14/sqlbuilder"
)

func TestMigrate(t *testing.T) {
	table := "t_user"
	updatedat.Field_UpdatedAt.Name = "Fauto_update_time"
	ddls := updatedat.Migrate(table, sqlbuilder.Dialect_mysql, sqlbuilder.SCENE_DDL_APPEND)
	fmt.Println(ddls)

}
