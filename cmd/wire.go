// This file uses wire to build all the depdendancies required

// +build wireinject


package cmd


import (
	"github.com/google/wire"
	config "github.com/spf13/viper"


	"github.com/snowzach/gorestapi/gorestapi"
	"github.com/snowzach/gorestapi/server"
	"github.com/snowzach/gorestapi/store/postgres"
)

// Create a new server
func NewServer() (*server.Server, error) {
    wire.Build(server.New, NewThingStore)
    return &server.Server{}, nil
}

// Create a new thing store
func NewThingStore() gorestapi.ThingStore {
	var thingStore gorestapi.ThingStore
	var err error
	switch config.GetString("storage.type") {
	case "postgres":
		thingStore, err = postgres.New()
	}
	if err != nil {
		logger.Fatalw("Database Error", "error", err)
	}
	return thingStore
}
