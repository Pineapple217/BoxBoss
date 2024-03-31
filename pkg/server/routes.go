package server

import (
	"github.com/Pineapple217/BoxBoss/pkg/handler"
)

func (s *Server) RegisterRoutes() {
	e := s.e

	e.Static("/static", "static/public")

	h := e.Group("/h")
	h.GET("/containers", handler.Containers)
	h.GET("/repo/createform", handler.RepoCreateForm)

	// api := e.Group("/api")

	e.GET("/repos", handler.Repos)
	// e.GET("/repo/create", handler.RepoCreate)
	e.POST("/repo/:id/build", handler.RepoBuild)
	e.POST("/repo/:id/update", handler.RepoUpdate)
	e.GET("/building", handler.Building)
	e.GET("/building_sse", handler.BuildingSSE)

	e.GET("/", handler.Home)
	// e.GET("/test", handler.Test)
}
