package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.com/gitedulab/learning-bot/modules/settings"
)

var engine *xorm.Engine

func init() {
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

}
