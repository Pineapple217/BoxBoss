package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Pineapple217/harbor-hawk/database"
	"github.com/Pineapple217/harbor-hawk/docker"
	"github.com/Pineapple217/harbor-hawk/handler"
	"github.com/Pineapple217/harbor-hawk/queue"
	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(echoMw.RequestLoggerWithConfig(echoMw.RequestLoggerConfig{
		LogStatus:  true,
		LogURI:     true,
		LogMethod:  true,
		LogLatency: true,
		LogValuesFunc: func(c echo.Context, v echoMw.RequestLoggerValues) error {
			slog.Info("request",
				"method", v.Method,
				"status", v.Status,
				"latency", v.Latency,
				"path", v.URI,
			)
			return nil

		},
	}))
	docker.Init()
	database.Init("file:database.db")

	queue.InitBuildQueue()

	e.Static("/static", "static/public")

	h := e.Group("/h")
	h.GET("/containers", handler.Containers)

	api := e.Group("/api")
	api.GET("/container/:id/update", handler.UpdateContainer)

	e.GET("/repos", handler.Repos)
	e.POST("/repo/:id/build", handler.RepoBuild)
	e.POST("/repo/:id/update", handler.RepoUpdate)
	e.GET("/building", handler.Building)
	e.GET("/building_sse", handler.BuildingSSE)

	e.GET("/", handler.Home)
	// e.GET("/test", handler.Test)

	go func() {
		if err := e.Start(":3000"); err != nil && err != http.ErrServerClosed {
			slog.Error("Shutting down the server", "error", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		slog.Error(err.Error())
	}
	// TODO: build cache

	// TODO: docker image
	// TODO: build and update button
	// TODO: private repos

	// TODO: more logging with slog
	// TODO: saving build logs

	// TODO: sh cant kill process warning
}
