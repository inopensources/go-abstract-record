package helpers

import "github.com/kataras/iris"

func Slicefy(ctx iris.Context, keys ...string) (keyValue []interface{}) {
	for _, key := range keys {
		keyValue = append(keyValue, key)
		keyValue = append(keyValue, ctx.Params().Get(key))
	}

	return keyValue
}