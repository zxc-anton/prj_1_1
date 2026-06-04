package middlewares

import (
	"server/pkg/token"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

func AuthMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		accessToken := ctx.Request.Header.Cookie("access_token")
		if len(accessToken) == 0 {
			ctx.Error("not valid token", fasthttp.StatusUnauthorized)
			return
		}
		id, err := token.DecodeToken(string(accessToken))
		if err != nil {
			ctx.Error("not valid token", fasthttp.StatusUnauthorized)
			log.Error().Msg(err.Error())
			return
		}

		ctx.SetUserValue("id", id)
		next(ctx)
	}
}
