package echo_swagger

import (
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/fatih/color"
	"github.com/labstack/echo"
	"github.com/ruanlianjun/gutils"
	"github.com/swaggo/swag"
)

const (
	defaultFilename     = "swagger.json"
	defaultMainFileName = "main.go"
)

type swagger struct {
	rootPath     string
	parser       *swag.Parser
	filename     string
	mainFileName string
	refresh      bool
	lock         sync.Mutex
}

func newSwagger() *swagger {
	rootPath, _ := os.Getwd()
	return &swagger{
		rootPath:     rootPath,
		parser:       swag.New(),
		filename:     defaultFilename,
		mainFileName: defaultMainFileName,
		refresh:      true,
		lock:         sync.Mutex{},
	}
}

type SwagOptions func(*swagger)

func WithSwaggerFilename(name string) SwagOptions {
	return func(s *swagger) {
		s.filename = name
	}
}

func WithRootPath(path string) SwagOptions {
	return func(s *swagger) {
		s.rootPath = path
	}
}

func WithSwagParse(parser *swag.Parser) SwagOptions {
	return func(s *swagger) {
		s.parser = parser
	}
}

func WithMainFilename(mainFilename string) SwagOptions {
	return func(s *swagger) {
		s.mainFileName = mainFilename
	}
}

func Refresh(refresh bool) SwagOptions {
	return func(s *swagger) {
		s.refresh = refresh
	}
}

var _swag *swagger = newSwagger()

func Swagger(options ...SwagOptions) echo.MiddlewareFunc {
	for _, item := range options {
		item(_swag)
	}
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			gutils.Recover()
			if _swag.refresh && _swag.lock.TryLock() {
				if err := generate(_swag); err != nil {
					return err
				}

			}
			return handlerFunc(context)
		}
	}
}

func Start(e *echo.Echo, addr string, showRoutes ...bool) error {

	e.GET("/swagger", func(context echo.Context) error {
		fs := filepath.Base(_swag.filename)
		return uiRender(context.Response(), filepath.Join("/swag/", fs))
	})

	dir := filepath.Join(_swag.rootPath, _swag.filename)

	if err := gutils.MkdirAll(dir); err != nil {
		_, _ = color.New(color.FgRed).Println(err)
		return err
	}

	// 静态文件代理

	fileDir := filepath.Dir(_swag.filename)
	fs := http.FileServer(http.Dir(filepath.Join(_swag.rootPath, fileDir)))

	e.GET("/swag/*", echo.WrapHandler(http.StripPrefix("/swag/", fs)))

	if len(showRoutes) > 0 && showRoutes[0] {
		displayRoutes(e)
	}

	svc := &http.Server{
		Handler: e,
		Addr:    addr,
	}

	return svc.ListenAndServe()
}

func displayRoutes(e *echo.Echo) {
	rows := make([][]any, 0, len(e.Routes()))
	for _, item := range e.Routes() {
		row := make([]any, 0, 3)
		row = append(row, color.New(color.FgYellow).Sprint(item.Name))
		row = append(row, color.New(color.FgGreen).Sprint(item.Method))
		row = append(row, color.New(color.FgBlue).Sprint(item.Path))
		rows = append(rows, row)
	}

	gutils.TerminalRender(
		gutils.SetRenderIsTerminal(),
		gutils.SetTableHeaders([]any{"名称", "方法", "URl"}),
		gutils.SetTableRows(rows),
	)
}
