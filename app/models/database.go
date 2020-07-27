package models

import (
	"github.com/Kamva/mgm"
	"github.com/revel/revel"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() {
	// Setup mgm default config
	err := mgm.SetDefaultConfig(nil, "viltasks", options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {

		revel.AppLog.Fatal("Database error", err)
	}
}
