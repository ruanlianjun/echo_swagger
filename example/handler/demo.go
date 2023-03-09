package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

// Hello
// @Summary 获取单个文章
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} string "成功"
// @Failure 400 {object} string "请求错误"
// @Failure 500 {object} string "内部错误"
// @Router /demo [get]
func Hello(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "hello")
}
