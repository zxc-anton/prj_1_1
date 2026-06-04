package http

import (
	"server/pkg/token"
	"server/pkg/types"
	"strconv"

	user_types "server/internal/service/user/types"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

func (h *HttpHandler) GetCurrentUser(ctx *fasthttp.RequestCtx) {
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

	user, err := h.user.GetCurrentUser(ctx, id)
	if err != nil {
		log.Error().Msg(err.Error())
		ctx.Error("not found", fasthttp.StatusNotFound)
		return
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.SetStatusCode(fasthttp.StatusOK)

	body, err := user.MarshalJSON()
	if err != nil {
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		log.Error().Msg(err.Error())
		return
	}

	ctx.Write(body)
}

func (h *HttpHandler) GetMyPlaylists(ctx *fasthttp.RequestCtx) {
	id, ok := ctx.UserValue("id").(int)

	if !ok {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	limit, offset := types.ParsePaginateParams(string(ctx.FormValue("limit")), string(ctx.FormValue("offset")))

	playlists, err := h.user.GetMyPlaylists(ctx, id, limit, offset)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		log.Error().Msg(err.Error())
		return
	}

	body, err := playlists.MarshalJSON()

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		log.Error().Msg(err.Error())
		return
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.SetStatusCode(fasthttp.StatusOK)

	ctx.Write(body)
}

func (h *HttpHandler) GetMyListenHistory(ctx *fasthttp.RequestCtx) {

	id, ok := ctx.UserValue("id").(int)

	if !ok {
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		return
	}

	limit, offset := types.ParsePaginateParams(string(ctx.FormValue("")), string(ctx.FormValue("")))

	tracks, err := h.user.GetMyListenHistory(ctx, id, limit, offset)

	if err != nil {
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		return
	}

	body, err := tracks.MarshalJSON()
	if err != nil {
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.Write(body)
}

func (h *HttpHandler) AddInFavorits(ctx *fasthttp.RequestCtx) {
	id, ok := ctx.UserValue("id").(int)

	if !ok {
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		return
	}

	var trackId types.Id
	if err := trackId.UnmarshalJSON(ctx.Request.Body()); err != nil {
		ctx.Error("not valid body", fasthttp.StatusUnprocessableEntity)
		log.Error().Msg(err.Error())
		return
	}

	if err := h.user.AddInFavorits(ctx, id, trackId.Id); err != nil {
		ctx.Error("not add", fasthttp.StatusBadRequest)
		log.Error().Msg(err.Error())
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
}

func (h *HttpHandler) CreatePlaylist(ctx *fasthttp.RequestCtx) {
	id, ok := ctx.UserValue("id").(int)

	if !ok {
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		return
	}

	var playlist user_types.CreatePlaylist
	playlist.UnmarshalJSON(ctx.Request.Body())

	if err := validate.Struct(playlist); err != nil {
		ctx.Error("not valid body", fasthttp.StatusUnprocessableEntity)
		log.Error().Msg(err.Error())
		return
	}

	if err := h.user.CreatePlaylist(ctx, id, playlist); err != nil {
		ctx.Error("bad request", fasthttp.StatusBadRequest)
		log.Error().Msg(err.Error())
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
}

func (h *HttpHandler) DeletePlaylist(ctx *fasthttp.RequestCtx) {
	id, ok := ctx.UserValue("id").(int)

	if !ok {
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		return
	}

	playlist_id_str, ok := ctx.UserValue("playlist_id").(string)
	playlist_id, err := strconv.Atoi(playlist_id_str)
	if !ok || err != nil {
		ctx.Error("no valid path", fasthttp.StatusUnprocessableEntity)
		return
	}

	if err := h.user.DeletePLaylist(ctx, id, playlist_id); err != nil {
		log.Error().Msg(err.Error())
	}

	ctx.SetStatusCode(fasthttp.StatusNoContent)
}

func (h *HttpHandler) DeleteTrackFromFavoris(ctx *fasthttp.RequestCtx) {
	id, ok := ctx.UserValue("id").(int)
	if !ok {
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		return
	}
	var trackId types.Id
	if err := trackId.UnmarshalJSON(ctx.Request.Body()); err != nil {
		ctx.Error("not valid body", fasthttp.StatusUnprocessableEntity)
		log.Error().Msg(err.Error())
		return
	}
	if err := h.user.DeleteTrackFromFavorits(ctx, id, trackId.Id); err != nil {
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		log.Error().Msg(err.Error())
		return
	}

	ctx.SetStatusCode(fasthttp.StatusNoContent)

}
