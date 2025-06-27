package mycontext

import (
	"context"
	"time"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
)

var ctx context.Context

func init() {
	// Mandatory to call config.Init() before using the context
	err := config.Init(time.Second * 10) // Load config file (yaml)
	if err != nil {
		log.DefaultLogger().Panicf(i18n.Translate(context.Background(), "Error while initializing config, err: %v"), err)
		panic(err)
	}

	ctx, err = config.TODOContext() // Global context
	if err != nil {
		log.DefaultLogger().Panicf(i18n.Translate(context.Background(), "Failed to create context: %v"), err)
	}
	log.DefaultLogger().Info(i18n.Translate(ctx, "Context initialized successfully"))
}

func GetContext() context.Context {
	return ctx
}
