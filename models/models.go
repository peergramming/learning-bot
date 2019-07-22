package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // MySQL driver support
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3" // SQLite driver support
	"gitlab.com/gitedulab/learning-bot/modules/settings"
	"log"
	"os"
	"xorm.io/core"
)

var (
	engine *xorm.Engine
	tables []interface{}
)

func init() {
	tables = append(tables,
		new(Repository),
		new(Report),
		new(Issue),
	)
}

// SetupEngine sets up the xorm engine according to the database configuration,
// and syncs the schema.
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
	case settings.PostgreSQL:
		dbAddr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
			dbConf.User, os.Getenv("POSTGRES_PASSWORD"),
			dbConf.Host, dbConf.Name, dbConf.SSLMode)
		engine, err = xorm.NewEngine("postgres", dbAddr)
	}

	if err != nil {
		log.Fatal("Unable to load database! ", err)
	}

	engine.TZLocation = settings.Config.Timezone
	engine.SetMapper(core.GonicMapper{})

	engine.Sync(tables...)
	return engine
}

func SetupTestEngine() {
	var err error
	testDBPath := "test.db"
	os.Remove(testDBPath)
	engine, err = xorm.NewEngine("sqlite3", testDBPath)
	if err != nil {
		log.Fatal("Unable to create test database! ", err)
	}
	engine.SetMapper(core.GonicMapper{})
	engine.Sync(tables...)
}
