package main

import (
	"flag"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var (
	outPath = flag.String("out", "", "out path for models")
	dsn     = flag.String("dsn", "host=localhost user=biz password=biz dbname=partner_service port=5433 sslmode=disable TimeZone=Europe/Moscow", "postgres dsn")
)

var dataMap = map[string]func(gorm.ColumnType) (dataType string){
	"uuid": func(columnType gorm.ColumnType) (dataType string) {
		return "uuid7.UUID"
	},
}

func main() {
	flag.Parse()

	g := gen.NewGenerator(gen.Config{
		OutPath:      "testModel/",
		ModelPkgPath: *outPath,
		Mode:         gen.WithoutContext,

		FieldNullable: true,
	})

	gormdb, _ := gorm.Open(postgres.Open(*dsn), &gorm.Config{})

	g.UseDB(gormdb)
	g.WithDataTypeMap(dataMap)
	g.ApplyBasic(
		g.GenerateAllTable(gen.FieldModify(fieldDBTag))...,
	)
	g.Execute()

	os.RemoveAll("testModel/")
}

func fieldDBTag(m gen.Field) gen.Field {
	value := m.ColumnName

	if _, ok := m.GORMTag["primaryKey"]; ok {
		value += ",primaryKey"
	}

	m.Tag.Set("db", value)
	m.Tag.Remove("json")

	m.GORMTag = make(map[string][]string)

	return m
}
