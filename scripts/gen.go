package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:       "./internal/data/layout/gen",
		Mode:          gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
		FieldSignable: true,
	})
	gormdb, err := gorm.Open(mysql.Open("root:Urie_308@tcp(127.0.0.1:3306)/xs"))
	if err != nil {
		panic(err)
	}
	g.UseDB(gormdb)
	dataTypeMap := map[string]func(columnType gorm.ColumnType) (dataType string){
		"decimal": func(detailType gorm.ColumnType) (dataType string) {
			return "decimal.Decimal"
		},
	}
	g.WithDataTypeMap(dataTypeMap)
	g.ApplyBasic(
	// 这里填写表名
	// g.GenerateModel("xs_user_id_pool"),
	// g.GenerateModel("xs_user_profile"),
	// g.GenerateModel("xs_user_mobile"),
	// g.GenerateModel("xs_user_platform"),
	)

	g.Execute()
}
