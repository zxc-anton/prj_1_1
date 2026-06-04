package run

import (
	"context"
	"fmt"
	"server/internal/config"
	"server/internal/handlers/http"
	"server/internal/service/artist"
	"server/internal/service/auth"
	"server/internal/service/playlist"
	"server/internal/service/search"
	"server/internal/service/track"
	"server/internal/service/user"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

func mustConnect(ctx context.Context, dsn string, minConns int, maxConns int) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}
	config.MinConns = int32(minConns)
	config.MaxConns = int32(maxConns)

	pool, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		panic(err)
	}

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}

	return pool
}

func Run(ctx context.Context, config config.Config) error {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DBName,
	)
	pool := mustConnect(ctx, dsn, config.Database.MinConns, config.Database.MaxConns)

	authService := auth.New(pool)
	userService := user.New(pool)
	trackService := track.New(pool)
	artistService := artist.New(pool)
	playlistService := playlist.New(pool)
	searchService := search.New(elasticsearch.WithAddresses("http://localhost:9200"))

	// IndexTracks(pool, searchService)

	log.Info().Msg("start server")

	return fasthttp.ListenAndServe(fmt.Sprintf("localhost:%d", config.Server.Port), http.NewHttpHandler(authService, userService, trackService, searchService, artistService, playlistService))
}
