package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"

	"example/handler"
	"github.com/ruanlianjun/echo_swagger"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func main() {
	e := echo.New()
	e.Pre(echo_swagger.Swagger(
		echo_swagger.WithSwaggerFilename("swagger/swagger.json"),
	))
	
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/demo", handler.Hello)

	e.HTTPErrorHandler = func(err error, context echo.Context) {
		fmt.Fprintf(os.Stdout, "URl:%s Method:%s Message:%s\n", context.Request().URL, context.Request().Method, err.Error())
	}

	e.Logger.Fatal(echo_swagger.Start(e, ":1323", true))
}
