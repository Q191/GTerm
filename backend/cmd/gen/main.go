package main

import (
	"github.com/OpenToolkitLab/GTerm/backend/dal/model"
	"gorm.io/gen"
	"os"
)

func models() []any {
	return []any{
		model.Host{},
		model.Credential{},
		model.Group{},
	}
}

func main() {
	_, err := os.Stat("main.go")
	if err != nil {
		panic("请在项目根目录执行 go run cmd/gen/mysql/main.go")
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "./backend/dal/query",
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.ApplyBasic(models()...)
	g.Execute()
}
