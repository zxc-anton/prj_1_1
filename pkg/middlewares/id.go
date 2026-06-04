package middlewares

import (
	"strconv"

	"github.com/valyala/fasthttp"
)

func IdMiddleware[T any](next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		idStr := ctx.UserValue("id").(string)
		id, err := strconv.Atoi(idStr)
		if err != nil {

		}
		ctx.SetUserValue("id", id)
		next(ctx)
	}
}
