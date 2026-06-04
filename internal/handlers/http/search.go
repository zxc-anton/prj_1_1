package http

import (
	"encoding/json"
	"server/pkg/types"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

func (h *HttpHandler) Search(ctx *fasthttp.RequestCtx) {
	q := ctx.FormValue("q")

	if len(q) == 0 {
		ctx.Error("not valid q params", fasthttp.StatusUnprocessableEntity)
		return
	}

	limit, offset := types.ParsePaginateParams(string(ctx.FormValue("limit")), string(ctx.FormValue("offset")))

	ids, err := h.search.SearchTrack(ctx, string(q), offset, limit)

	if err != nil {
		log.Error().Err(err).Send()
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

	ctx.Response.Header.SetContentType("application/json")
	ctx.Write(body)
}
