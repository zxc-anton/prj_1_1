package http

import (
	"encoding/json"
	"errors"
	"server/internal/service/artist"
	"server/pkg/types"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

func (h *HttpHandler) GetArtistById(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		log.Error().Err(err).Send()
		ctx.Error("not valid id", fasthttp.StatusUnprocessableEntity)
		return
	}

	artistResponse, err := h.artist.GetArtistById(ctx, id)

	if err != nil {
		if errors.Is(err, artist.ErrNotFound) {
			ctx.Error("not found", fasthttp.StatusNotFound)
			return
		}
		log.Error().Err(err).Send()
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		return
	}

	body, err := artistResponse.MarshalJSON()

	if err != nil {
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		log.Error().Err(err).Send()
		return
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.Write(body)
}

func (h *HttpHandler) GetArtistAlbums(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		log.Error().Err(err).Send()
		ctx.Error("not valid id", fasthttp.StatusUnprocessableEntity)
		return
	}

	limit, offset := types.ParsePaginateParams(string(ctx.FormValue("limit")), string(ctx.FormValue("offset")))

	albums, err := h.artist.GetArtistAlbums(ctx, id, limit, offset)

	if err != nil {
		log.Error().Err(err).Send()
		ctx.Error("intrnal error", fasthttp.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(albums)

	if err != nil {
		log.Error().Err(err).Send()
		ctx.Error("intrnal error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.Write(body)
}

func (h *HttpHandler) GetArtistPopularTracks(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error().Err(err).Send()
		ctx.Error("not valid id", fasthttp.StatusUnprocessableEntity)
		return
	}

	limit, offset := types.ParsePaginateParams(string(ctx.FormValue("limit")), string(ctx.FormValue("offset")))
	ids, err := h.artist.GetArtistPopularTracksIDS(ctx, id, limit, offset)

	if err != nil {
		log.Error().Err(err).Send()
		ctx.Error("intrnal error", fasthttp.StatusInternalServerError)
		return
	}

	tracks, err := h.track.GetTracksByIDS(ctx, ids)
	if err != nil {
		log.Error().Err(err).Send()
		ctx.Error("intrnal error", fasthttp.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(tracks)
	if err != nil {
		log.Error().Err(err).Send()
		ctx.Error("intrnal error", fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.Write(body)
}
