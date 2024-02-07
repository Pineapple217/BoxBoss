package handler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Pineapple217/harbor-hawk/view"
	"github.com/labstack/echo/v4"
)

func BuildUI(c echo.Context) error {
	return render(c, view.Build())
}

func BuildSSE(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")

	for i := 1; i <= 3; i++ {

		// c.String(http.StatusOK, fmt.Sprintf("Iteration %d\n", i+1))
		// c.Request().Write()
		fmt.Fprint(c.Response(), buildSSE("message", strconv.Itoa(i)))
		c.Response().Flush()
		time.Sleep(3 * time.Second)
	}
	fmt.Fprint(c.Response(), buildSSE("message", "<div id='sse-feed' hx-swap-oob='true'></div>"))
	c.Response().Flush()

	return nil
}
