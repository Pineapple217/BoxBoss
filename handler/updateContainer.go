package handler

import (
	"fmt"
	"net/http"

	"github.com/Pineapple217/harbor-hawk/docker"
	"github.com/labstack/echo/v4"
)

func UpdateContainer(c echo.Context) error {
	id := c.Param("id")
	fmt.Println(id)

	err := docker.UpdateContainer(id)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "aaa")
}

func Build(c echo.Context) error {
	// start := time.Now()
	// err := docker.BuildAndUploadImage("https://github.com/Pineapple217/cicd-testing", "user", "pwd")
	// elapsed := time.Since(start)
	// log.Printf("Building took %s", elapsed)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }

	return c.String(http.StatusOK, "bbb")
}
