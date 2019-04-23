package controllers

import (
	"github.com/kataras/iris"
)

func CheckReturn(ctx iris.Context, data interface{}, err error, message string) {
	if err == nil {
		ctx.JSON(Response{200, data, message})
	} else {
		ctx.JSON(Response{500, err, "Ocorreu algum erro dentro do sistema!!"})
	}
}