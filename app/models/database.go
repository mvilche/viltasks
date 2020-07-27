package models

import (
	"os"

	"github.com/Kamva/mgm"
	"github.com/revel/config"
	"github.com/revel/revel"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() {
	// Setup mgm default config

	c, err := config.ReadDefault("./conf/app.conf")
	if err != nil {
		revel.AppLog.Error("", err)
		os.Exit(1)
	}

	var host, _ = c.String("database", "database.host")
	var port, _ = c.String("database", "database.port")
	var user, _ = c.String("database", "database.username")
	var name, _ = c.String("database", "database.name")
	var pass, _ = c.String("database", "database.password")

	if pass == "" || user == "" {

		err2 := mgm.SetDefaultConfig(nil, name, options.Client().ApplyURI("mongodb://"+host+":"+port+""))
		if err2 != nil {
			revel.AppLog.Fatal("Database error", err2)
		}
	} else {

		err2 := mgm.SetDefaultConfig(nil, name, options.Client().ApplyURI("mongodb://"+user+":"+pass+"@"+host+":"+port+""))
		if err2 != nil {
			revel.AppLog.Fatal("Database error", err2)
		}

	}

}
