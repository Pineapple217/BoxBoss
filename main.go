package main

import (
	"fmt"

	"github.com/Pineapple217/harbor-hawk/database"
	"github.com/Pineapple217/harbor-hawk/docker"
	"github.com/Pineapple217/harbor-hawk/handler"
	"github.com/Pineapple217/harbor-hawk/queue"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			fmt.Printf("REQUEST: uri: %v, status: %v\n", v.URI, v.Status)
			return nil
		},
	}))
	docker.Init()
	database.Init("file:database.db")

	queue.InitBuildQueue()

	e.Static("/static", "static/public")

	h := e.Group("/h")
	h.GET("/containers", handler.Containers)
	h.GET("/building_sse", handler.BuildingSSE)

	api := e.Group("/api")
	api.GET("/container/:id/update", handler.UpdateContainer)
	e.GET("/build", handler.Build)

	e.GET("/repos", handler.Repos)
	e.POST("/repo/:id/build", handler.RepoBuild)
	e.GET("/building", handler.Building)

	e.GET("/", handler.Home)
	e.GET("build", handler.BuildUI)
	e.GET("build_sse", handler.BuildSSE)
	// e.GET("/test", handler.Test)
	e.Start(":3000")
}
