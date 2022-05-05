package main

import (
	"flag"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"kratosx-fashion/app/system/internal/data/model"
)

var outPath = flag.String("output", "./internal/data/linq", "output path")
var dsn = flag.String("dsn", "root:123456@(127.0.0.1:3306)/violet_test?charset=utf8mb4&parseTime=True&loc=Local", "dsn")

//go:generate
//go:generate go run generate.go -output=./internal/data/linq
func main() {
	// specify the output directory (default: "./query")
	// ### if you want to query without context constrain, set mode gen.WithoutContext ###
	g := gen.NewGenerator(gen.Config{
		OutPath: *outPath,
		/* Mode: gen.WithoutContext|gen.WithDefaultQuery*/
		//if you want the nullable field generation property to be pointer type, set FieldNullable true
		/* FieldNullable: true,*/
		//if you want to assign field which has default value in `Create` API, set FieldCoverable true, reference: https://gorm.io/docs/create.html#Default-Values
		FieldCoverable: true,
		//if you want to generate index tags from database, set FieldWithIndexTag true
		FieldWithIndexTag: true,
		//if you want to generate type tags from database, set FieldWithTypeTag true
		/* FieldWithTypeTag: true,*/
		//if you need unit tests for query code, set WithUnitTest true
		/* WithUnitTest: true, */
	})
	// reuse the database connection in Project or create a connection here
	// if you want to use GenerateModel/GenerateModelAs, UseDB is necessary or it will panic
	db, _ := gorm.Open(mysql.Open(*dsn))
	g.UseDB(db)

	// apply basic crud api on structs or table models which is specified by table name with function
	// GenerateModel/GenerateModelAs. And generator will generate table models' code when calling Excute.
	g.ApplyBasic(
		model.UserLog{},
		model.User{}, model.UserRole{},
		model.Role{}, model.RoleResource{},
		model.ResourceMenu{},
	)

	// execute the action of code generation
	g.Execute()
}
