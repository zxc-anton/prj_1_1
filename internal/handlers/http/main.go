package http

import (
	"server/internal/service/album"
	"server/internal/service/artist"
	"server/internal/service/auth"
	"server/internal/service/playlist"
	"server/internal/service/search"
	"server/internal/service/track"
	"server/internal/service/user"
	"server/pkg/middlewares"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/pprofhandler"
)

type HttpHandler struct {
	auth     *auth.AuthService
	user     *user.UserService
	track    *track.TrackService
	search   *search.SearchService
	artist   *artist.ArtistService
	playlist *playlist.PLaylistService
	album    *album.AlbumService
}

func NewHttpHandler(
	auth *auth.AuthService,
	user *user.UserService,
	track *track.TrackService,
	search *search.SearchService,
	artist *artist.ArtistService,
	playlist *playlist.PLaylistService,
	album *album.AlbumService,
) fasthttp.RequestHandler {
	h := HttpHandler{
		auth:     auth,
		user:     user,
		track:    track,
		search:   search,
		artist:   artist,
		playlist: playlist,
		album:    album,
	}

	r := router.New()

	r.GET("/debug/pprof/{profile:*}", pprofhandler.PprofHandler)
	r.POST("/auth/users", h.RequestUser)
	r.POST("/auth/login", h.LoginUser)
	r.POST("/auth/refresh_token", h.RefreshToken)
	r.GET("/auth/current_user", h.GetCurrentUser)
	r.GET("/me/playlists", middlewares.AuthMiddleware(h.GetMyPlaylists))
	r.GET("/me/listen_history", middlewares.AuthMiddleware(h.GetMyListenHistory))
	r.POST("/me/favortis", middlewares.AuthMiddleware(h.AddInFavorits))
	r.POST("/me/playlists", middlewares.AuthMiddleware(h.CreatePlaylist))
	r.GET("/me/favortis", middlewares.AuthMiddleware(h.GetFavorits))
	r.DELETE("/me/playlists/{playlist_id}", middlewares.AuthMiddleware(h.DeletePlaylist))
	r.DELETE("/me/favortis", middlewares.AuthMiddleware(h.DeleteTrackFromFavoris))
	r.GET("/health", HealthCheck)
	r.GET("/track/popular", h.GetPopularTracks)
	r.GET("/track/{id}", h.GetTrackById)
	r.GET("/track/{id}/authors", h.GetTrackAuthors)
	r.GET("/search", h.Search)
	r.GET("/artist/{id}", h.GetArtistById)
	r.GET("/artist/{id}/albums", h.GetArtistAlbums)
	r.GET("/artist/{id}/tracks", h.GetArtistPopularTracks)
	r.GET("/playlist/{id}", h.GetPLaylistTracks)
	r.GET("/album/{id}", h.GetAlbumById)
	r.GET("/album/{id}/tracks", h.GetAlbumTracks)

	return middlewares.LoggMiddleware(middlewares.CORSMiddleware(r.Handler))
}

func HealthCheck(ctx *fasthttp.RequestCtx) {
	ctx.Write([]byte("ok"))
}
