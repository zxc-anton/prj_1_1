package http

import (
	"server/internal/service/auth/types"
	"time"

	"github.com/go-playground/validator"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (h HttpHandler) RequestUser(ctx *fasthttp.RequestCtx) {
	var user types.RequestUser

	if err := user.UnmarshalJSON(ctx.Request.Body()); err != nil {
		log.Error().Msg(err.Error())
		ctx.Error("not valid body", fasthttp.StatusUnprocessableEntity)
		return
	}

	if err := validate.Struct(user); err != nil {
		log.Error().Msg(err.Error())
		ctx.Error("not valid body", fasthttp.StatusUnprocessableEntity)
		return
	}

	if err := h.auth.RegisteUser(ctx, user); err != nil {
		log.Error().Msg(err.Error())
		ctx.Error("not valid body", fasthttp.StatusUnprocessableEntity)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
}

func (h HttpHandler) LoginUser(ctx *fasthttp.RequestCtx) {
	var params types.LoginUser

	if err := params.UnmarshalJSON(ctx.Request.Body()); err != nil {
		log.Error().Msg(err.Error())
		ctx.Error("not valid body", fasthttp.StatusUnprocessableEntity)
		return
	}

	if err := validate.Struct(params); err != nil {
		log.Error().Msg(err.Error())
		ctx.Error("not valid body", fasthttp.StatusUnprocessableEntity)
		return
	}

	accessToken, refreshToken, err := h.auth.LoginUser(ctx, params)
	if err != nil {
		log.Error().Msg(err.Error())
		ctx.Error("cant auntifacate", fasthttp.StatusUnauthorized)
		return
	}

	accessCookie := fasthttp.AcquireCookie()
	refreshCookie := fasthttp.AcquireCookie()

	defer fasthttp.ReleaseCookie(accessCookie)
	defer fasthttp.ReleaseCookie(refreshCookie)

	accessCookie.SetKey("access_token")
	accessCookie.SetValue(accessToken)
	accessCookie.SetHTTPOnly(true)
	accessCookie.SetPath("/")

	refreshCookie.SetKey("refresh_token")
	refreshCookie.SetValue(refreshToken)
	refreshCookie.SetHTTPOnly(true)
	refreshCookie.SetPath("/")

	ctx.Response.Header.SetCookie(accessCookie)
	ctx.Response.Header.SetCookie(refreshCookie)
}

func (h HttpHandler) Logout(ctx *fasthttp.RequestCtx) {

	accessCookie := fasthttp.AcquireCookie()
	refreshCookie := fasthttp.AcquireCookie()

	defer fasthttp.ReleaseCookie(accessCookie)
	defer fasthttp.ReleaseCookie(refreshCookie)

	accessCookie.SetKey("access_token")
	accessCookie.SetValue("")
	accessCookie.SetExpire(time.Now())
	accessCookie.SetMaxAge(-1)

	refreshCookie.SetKey("refresh_token")
	refreshCookie.SetValue("")
	refreshCookie.SetExpire(time.Now())
	refreshCookie.SetMaxAge(-1)

	ctx.Response.Header.SetCookie(accessCookie)
	ctx.Response.Header.SetCookie(refreshCookie)
}

func (h *HttpHandler) RefreshToken(ctx *fasthttp.RequestCtx) {
	refreshToken := string(ctx.Request.Header.Cookie("refresh_token"))

	if len(refreshToken) == 0 {
		ctx.Error("not valid token", fasthttp.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := h.auth.RefreshToken(ctx, refreshToken)

	if err != nil {
		log.Error().Err(err).Send()
		ctx.Error("not valid token", fasthttp.StatusUnauthorized)
		return
	}

	accessCookie := fasthttp.AcquireCookie()
	refreshCookie := fasthttp.AcquireCookie()

	defer fasthttp.ReleaseCookie(accessCookie)
	defer fasthttp.ReleaseCookie(refreshCookie)

	accessCookie.SetKey("access_token")
	accessCookie.SetValue(accessToken)
	accessCookie.SetHTTPOnly(true)
	accessCookie.SetPath("/")

	refreshCookie.SetKey("refresh_token")
	refreshCookie.SetValue(refreshToken)
	refreshCookie.SetHTTPOnly(true)
	refreshCookie.SetPath("/")

	ctx.Response.Header.SetCookie(accessCookie)
	ctx.Response.Header.SetCookie(refreshCookie)

}
