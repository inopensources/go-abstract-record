package helpers

import (
	"github.com/kataras/iris"
)

func Slicefy(ctx iris.Context, keys ...string) (keyValue []interface{}) {
	for _, key := range keys {
		keyValue = append(keyValue, key)
		keyValue = append(keyValue, ctx.Params().Get(key))
	}

	return keyValue
}

func GetUrlParams(ctx iris.Context) (urlParams map[string]string) {
	urlParams = ctx.URLParams()
	delete(urlParams, "offset")
	delete(urlParams, "pagesize")

	return urlParams
}

func GetPaginationValues(ctx iris.Context) (offset, pagesize int) {
	offset, _ = ctx.URLParamInt("offset")
	pagesize, _ = ctx.URLParamInt("pagesize")

	return offset, pagesize
}