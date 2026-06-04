package http

import (
	"encoding/json"
	"server/pkg/types"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

func (h *HttpHandler) GetPLaylistTracks(ctx *fasthttp.RequestCtx) {
	idStr := ctx.UserValue("id").(string)
	id, err := strconv.Atoi(idStr)

	if err != nil {
		log.Error().Err(err).Send()
		ctx.Error("not valid id", fasthttp.StatusUnprocessableEntity)
		return
	}

	limit, offset := types.ParsePaginateParams(string(ctx.FormValue("limit")), string(ctx.FormValue("offset")))
	ids, err := h.playlist.GetPLaylistTracks(ctx, id, limit, offset)

	if err != nil {
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		log.Error().Err(err).Send()
		return
	}

	tracks, err := h.track.GetTracksByIDS(ctx, ids)
	if err != nil {
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		log.Error().Err(err).Send()
		return
	}

	body, err := json.Marshal(tracks)
	if err != nil {
		ctx.Error("internal error", fasthttp.StatusInternalServerError)
		log.Error().Err(err).Send()
		return
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.Write(body)
}
