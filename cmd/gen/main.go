//go:generate go run .
package main

import (
	"github.com/yikakia/nga_grep/model"
	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:           "../../model/gen",
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		Mode:              gen.WithoutContext | gen.WithQueryInterface | gen.WithDefaultQuery,
	})
	g.ApplyBasic(&model.ThreadLatestData{})
	g.ApplyBasic(&model.ThreadCount{})
	g.Execute()
}
