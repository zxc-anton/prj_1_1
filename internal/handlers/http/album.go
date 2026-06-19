package http

import (
	"encoding/json"
	"errors"
	"server/internal/service/track"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

func (h *HttpHandler) GetAlbumById(ctx *fasthttp.RequestCtx) {
	idStr, _ := ctx.UserValue("id").(string)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		log.Error().Err(err).Send()
		ctx.Error("not valid id", fasthttp.StatusUnprocessableEntity)
		return
	}

	album, err := h.album.GetAlbumById(ctx, id)
	if err != nil {
		log.Error().Err(err).Send()
		if errors.Is(err, track.ErrNotFound) {
			ctx.Error("not found", fasthttp.StatusNotFound)
			return
		}
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		return
	}

	body, err := album.MarshalJSON()

	if err != nil {
		log.Error().Err(err).Send()
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.Write(body)
	ctx.Response.Header.SetContentType("application/json")
}

func (h *HttpHandler) GetAlbumTracks(ctx *fasthttp.RequestCtx) {
	idStr, _ := ctx.UserValue("id").(string)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		log.Error().Err(err).Send()
		ctx.Error("not valid id", fasthttp.StatusUnprocessableEntity)
		return
	}

	ids, err := h.album.GetAlbumTracks(ctx, id)
	if err != nil {
		log.Error().Err(err).Send()
		if errors.Is(err, track.ErrNotFound) {
			ctx.Error("not found", fasthttp.StatusNotFound)
			return
		}
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		return
	}

	tracks, err := h.track.GetTracksByIDS(ctx, ids)
	if err != nil {
		log.Error().Err(err).Send()
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(tracks)
	if err != nil {
		log.Error().Err(err).Send()
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.Write(body)
	ctx.Response.Header.SetContentType("application/json")
}
