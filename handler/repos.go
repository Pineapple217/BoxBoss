package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Pineapple217/harbor-hawk/database"
	"github.com/Pineapple217/harbor-hawk/docker"
	"github.com/Pineapple217/harbor-hawk/queue"
	"github.com/Pineapple217/harbor-hawk/view"
	"github.com/labstack/echo/v4"
)

func Repos(c echo.Context) error {
	quaries := database.GetQueries()
	repos, err := quaries.ListRepos(c.Request().Context())
	if err != nil {
		return err
	}
	return render(c, view.Repos(repos))
}

func RepoBuild(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return err
	}
	queries := database.GetQueries()
	repo, err := queries.GetRepo(c.Request().Context(), id)
	if err != nil {
		return err
	}

	queue := queue.GetBuildQueue()
	queue.Enqueue(docker.BuildSettings{
		Repo: &repo,
	})

	c.Response().Header().Add("HX-Redirect", "/building")
	return c.NoContent(http.StatusAccepted)
}

func RepoUpdate(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return err
	}
	queries := database.GetQueries()
	repo, err := queries.GetRepo(c.Request().Context(), id)
	if err != nil {
		return err
	}
	composePath := "/opt/stacks/" + repo.ComposeFile.String + "/docker-compose.yml"
	docker.ComposePull(composePath, repo.ComposeService.String)
	docker.ComposeStop(composePath, repo.ComposeService.String)
	docker.ComposeRemove(composePath, repo.ComposeService.String)
	docker.ComposeUp(composePath, repo.ComposeService.String)

	c.Response().Header().Add("HX-Redirect", "/repos")
	return c.NoContent(http.StatusAccepted)
}

func Building(c echo.Context) error {
	return render(c, view.Building())
}

// func Test(c echo.Context) error {
// 	// docker.Test()
// 	// return c.String(200, "")
// 	return render(c, view.BuildingTest())
// }

// TODO: if more then 1 page open this breaks
func BuildingSSE(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")

	queue := queue.GetBuildQueue()

	ctx := c.Request().Context()
	// TODO: only flush every 500ms or smth
	for {
		select {
		case result := <-queue.BuildLogsChannel:
			fmt.Fprint(c.Response(), buildSSE("message", result))
			c.Response().Flush()
		case <-ctx.Done():
			return nil
		}
	}
}
