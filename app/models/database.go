package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/revel/config"
	"github.com/revel/revel"
	"os"
)

func OpenSQL() (*gorm.DB, error) {

	c, err := config.ReadDefault("./conf/app.conf")
	if err != nil {
		revel.AppLog.Error("", err)
		os.Exit(1)
	}

	var url, _ = c.String("database", "database.url")

	db, err := gorm.Open("sqlite3", url)
	if err != nil {

		return db, err
	}
	//defer db.Close()

	return db, err
}

func CloseSQL(db *gorm.DB) {

	db.Close()

}
