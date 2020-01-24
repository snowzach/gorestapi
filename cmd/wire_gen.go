// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package cmd

import (
	"github.com/snowzach/gorestapi/gorestapi"
	"github.com/snowzach/gorestapi/server"
	"github.com/snowzach/gorestapi/store/postgres"
	"github.com/spf13/viper"
)

import (
	_ "net/http/pprof"
)

// Injectors from wire.go:

func NewServer() (*server.Server, error) {
	serverServer, err := server.New()
	if err != nil {
		return nil, err
	}
	return serverServer, nil
}

// wire.go:

// Create a new thing store
func NewThingStore() gorestapi.ThingStore {
	var thingStore gorestapi.ThingStore
	var err error
	switch viper.GetString("storage.type") {
	case "postgres":
		thingStore, err = postgres.New()
	}
	if err != nil {
		logger.Fatalw("Database Error", "error", err)
	}
	return thingStore
}
