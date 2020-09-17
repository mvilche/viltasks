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
		revel.AppLog.Error("", err)
		os.Exit(1)
	}

	var url, _ = c.String("database", "database.url")

	// github.com/mattn/go-sqlite3
	db, err := gorm.Open(sqlite.Open(url), &gorm.Config{

		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {

		return db, err
	}
	//defer db.Close()
	return db, err
}
