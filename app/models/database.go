package models

import (
	"os"

	"github.com/revel/config"
	"github.com/revel/revel"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenSQL() (*gorm.DB, error) {

	c, err := config.ReadDefault("./conf/app.conf")
	if err != nil {
		revel.AppLog.Error("Error al obtener app.conf file", err)
		os.Exit(1)
	}

	var url, err2 = c.String("database", "database.url")

	if err2 != nil {
		revel.AppLog.Error("Error al obtener database url", err)
		os.Exit(1)
	}
	revel.AppLog.Debug("Se encontro database url: " + url)

	db, err := gorm.Open(sqlite.Open(url), &gorm.Config{

		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		revel.AppLog.Error("", err)
		return db, err
	}
	//defer db.Close()
	return db, err
}
