package initialize

import (
	"fmt"
	"github.com/OpenToolkitLab/GTerm/backend/consts"
	"github.com/OpenToolkitLab/GTerm/backend/dal/model"
	"github.com/OpenToolkitLab/GTerm/backend/dal/query"
	"github.com/OpenToolkitLab/GTerm/backend/pkg/storage"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	db    *gorm.DB
	Query *query.Query
}

func InitDatabase() *Database {
	database := &Database{}
	localStorage := storage.NewLocalStorage(fmt.Sprintf("%s.%s", consts.ProjectName, consts.DatabaseDriver))
	if err := localStorage.CreateDirectory(); err != nil {
		panic(err)
	}
	if !localStorage.DatabaseExist() {
		if err := database.connect(localStorage.Path); err != nil {
			panic(err)
		}
		if err := database.autoMigrate(); err != nil {
			panic(err)
		}
		return database
	}
	if err := database.connect(localStorage.Path); err != nil {
		panic(err)
	}
	return database
}

func (d *Database) connect(dsn string) error {
	database, err := gorm.Open(sqlite.Open(dsn))
	if err != nil {
		return err
	}
	d.db = database
	d.Query = query.Use(d.db)
	return nil
}

func (d *Database) autoMigrate() error {
	models := []any{
		model.Connection{},
		model.Credential{},
		model.ConnectionGroup{},
	}
	return d.db.AutoMigrate(models...)
}
