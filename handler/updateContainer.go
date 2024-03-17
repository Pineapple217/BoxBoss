package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UpdateContainer(c echo.Context) error {
	id := c.Param("id")
	fmt.Println(id)

	// err := docker.UpdateContainer(id)
	// if err != nil {
	// 	return err
	// }
	return c.String(http.StatusOK, "aaa")
}
