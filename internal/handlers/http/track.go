package http

import (
	"encoding/json"
	"errors"
	"server/internal/service/track"
	track_service "server/internal/service/track"
	"server/pkg/types"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"golang.org/x/sync/singleflight"
)

func (h *HttpHandler) GetTrackById(ctx *fasthttp.RequestCtx) {
	idStr, _ := ctx.UserValue("id").(string)

	id, err := strconv.Atoi(idStr)
	if err != nil {

	}

	if err != nil {
		log.Error().Err(err).Send()
		ctx.Error("not valid id", fasthttp.StatusUnprocessableEntity)
		return
	}

	track, err := h.track.GetTrackById(ctx, id)
	if err != nil {
		log.Error().Err(err).Send()
		if errors.Is(err, track_service.ErrNotFound) {
			ctx.Error("not found", fasthttp.StatusNotFound)
			return
		}
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		return
	}

	body, err := track.MarshalJSON()
	if err != nil {
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.Write(body)
}

func (h *HttpHandler) GetTrackAuthors(ctx *fasthttp.RequestCtx) {
	idStr, _ := ctx.UserValue("id").(string)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		log.Error().Err(err).Send()
		ctx.Error("not valid id", fasthttp.StatusUnprocessableEntity)
		return
	}

	authors, err := h.track.GetTrackAuthors(ctx, id)

	if err != nil {
		log.Error().Err(err).Send()
		ctx.Error("internal server error", fasthttp.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(authors)

	if err != nil {
		log.Error().Err(err).Send()
		ctx.Error("internal server error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.Write(body)
}

var TrackGroup singleflight.Group

func (h *HttpHandler) GetPopularTracks(ctx *fasthttp.RequestCtx) {
	limit, offset := types.ParsePaginateParams(string(ctx.FormValue("limit")), string(ctx.FormValue("offset")))

	tracksAny, err, _ := TrackGroup.Do("popular", func() (any, error) {
		tracks, err := h.track.GetPopularTracks(ctx, limit, offset)
		return tracks, err
	})

	tracks, _ := tracksAny.([]track.Track)

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

	ctx.Response.Header.SetContentType("application/json")
	ctx.Write(body)
}
