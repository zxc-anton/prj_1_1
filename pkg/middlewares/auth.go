package middlewares

import (
	"errors"
	"server/pkg/token"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

var ErrNotValidToken = errors.New("not valid token")

func AuthUser(accessToken string) (int, error) {
	if len(accessToken) == 0 {
		return 0, ErrNotValidToken
	}

	id, err := token.DecodeToken(accessToken)
	return id, err
}

func AuthMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		accessToken := ctx.Request.Header.Cookie("access_token")
		id, err := AuthUser(string(accessToken))
		if err != nil {
			ctx.Error("not valid token", fasthttp.StatusUnauthorized)
			log.Error().Msg(err.Error())
			return
		}

		ctx.SetUserValue("id", id)
		next(ctx)
	}
}

func MaybeAuthMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		accessToken := ctx.Request.Header.Cookie("access_token")
		id, err := AuthUser(string(accessToken))
		if err == nil {
			ctx.SetUserValue("user_id", id)
		}
		next(ctx)
	}
}
