package middlewares

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

func LoggMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		defer func() {
			if err := recover(); err != nil {
				ctx.SetStatusCode(500)
				log.Info().
					Int("status_code", ctx.Response.StatusCode()).
					Str("method", string(ctx.Method())).
					Dur("latency", time.Since(start)).
					Any("error", err).
					Msg("error")
			}
		}()

		next(ctx)

		log.Info().
			Int("status_code", ctx.Response.StatusCode()).
			Str("method", string(ctx.Method())).
			Str("path", string(ctx.Path())).
			Dur("latency", time.Since(start)).
			Msg("ok")
	}

}
