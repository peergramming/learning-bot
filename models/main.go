package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"xorm.io/core"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.com/gitedulab/learning-bot/modules/settings"
	"log"
	"os"
)

var (
	engine *xorm.Engine
	tables []interface{}
)

func SetupEngine() *xorm.Engine {
	var err error
	dbConf := &settings.Config.DatabaseConfiguration

	switch dbConf.Type {
	case settings.SQLite:
		engine, err = xorm.NewEngine("sqlite3", settings.Config.DatabaseConfiguration.Path)
	case settings.MySQL:
		dbAddr := fmt.Sprintf("%s:%s@%s/%s?charset=utf8&tls=%s",
			dbConf.User, os.Getenv("MYSQL_PASSWORD"),
			dbConf.Host, dbConf.Name, dbConf.SSLMode)
		engine, err = xorm.NewEngine("mysql", dbAddr)
	}

	if err != nil {
		log.Fatal("Unable to load database! ", err)
	}

	tables = append(tables,
		new(Repository),
		new(Report),
	)
	engine.SetMapper(core.GonicMapper{})

	engine.Sync(tables...)
	return engine
}
