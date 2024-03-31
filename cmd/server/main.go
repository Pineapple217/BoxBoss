package main

import (
	"log/slog"
	"os"
	"os/signal"

	"github.com/Pineapple217/BoxBoss/pkg/database"
	"github.com/Pineapple217/BoxBoss/pkg/docker"
	"github.com/Pineapple217/BoxBoss/pkg/queue"
	"github.com/Pineapple217/BoxBoss/pkg/server"
)

func main() {
	server := server.NewServer()
	server.RegisterRoutes()
	server.ApplyMiddleware()

	docker.Init()
	database.Init("file:data/database.db")

	queue.InitBuildQueue()

	server.Start()

	// TODO: build cache

	// TODO: docker image
	// TODO: build and update button
	// TODO: private repos

	// TODO: more logging with slog
	// TODO: saving build logs

	// TODO: sh cant kill process warning

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	slog.Info("Received an interrupt signal, exiting...")

	server.Stop()
}
