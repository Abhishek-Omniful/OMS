package main

import (
	"log"

	initpkg "github.com/Abhishek-Omniful/OMS/init"
	"github.com/Abhishek-Omniful/OMS/mycontext"

	middlewares "github.com/Abhishek-Omniful/OMS/pkg/middleware"
	"github.com/Abhishek-Omniful/OMS/pkg/router"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/http"
)



func main() {
	ctx := mycontext.GetContext()
	initpkg.Initialize()
	server := http.InitializeServer(
		config.GetString(ctx, "server.port"),            // Port to listen
		config.GetDuration(ctx, "server.read_timeout"),  // Read timeout
		config.GetDuration(ctx, "server.write_timeout"), // Write timeout
		config.GetDuration(ctx, "server.idle_timeout"),  // Idle timeout
		false,
	)
	server.Use(middlewares.LogRequest(ctx)) // Middleware for logging requests
	router.Initialize(server)
	err := server.StartServer(config.GetString(ctx, "server.name"))
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}

}
