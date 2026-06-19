package middlewares

import "github.com/valyala/fasthttp"

func CORSMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		next(ctx)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "http://localhost:5173")
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	}
}
