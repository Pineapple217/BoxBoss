package handler

import (
	"github.com/Pineapple217/BoxBoss/pkg/docker"
	"github.com/Pineapple217/BoxBoss/pkg/view"
	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {
	// containers := docker.Ps()
	// return render(c, view.Base(docker.Ps()))
	return render(c, view.ContainersBase())
}

func Containers(c echo.Context) error {
	return render(c, view.Containers(docker.Ps()))
}
