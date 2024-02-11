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

// var mainCh chan string

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
	// mainCh = make(chan string, 50)
	// start := time.Now()

	// go docker.BuildAndUploadImage(repo, "user", "pwd", mainCh)
	queue := queue.GetBuildQueue()
	queue.Enqueue(docker.BuildSettings{
		Repo: &repo,
	})

	// go func() {
	// 	for s := range ch {
	// 		fmt.Println("Received:", s)
	// 	}
	// }()

	// elapsed := time.Since(start)
	// log.Printf("Building took %s", elapsed)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }
	// for result := range mainCh {
	// 	println(result)
	// }
	// wg.Wait()
	c.Response().Header().Add("HX-Redirect", "/building")
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
	for {
		select {
		case result := <-queue.BuildLogsChannel:
			fmt.Fprint(c.Response(), buildSSE("message_encoded", result))
			c.Response().Flush()
		case <-ctx.Done():
			return c.NoContent(204)
		}
	}
}
