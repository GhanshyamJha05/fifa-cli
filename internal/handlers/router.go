package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/GhanshyamJha05/fifa-cli/internal/middleware"
	"github.com/GhanshyamJha05/fifa-cli/internal/service"
)

// NewRouter builds the HTTP mux with all routes and middleware.
func NewRouter(svc *service.Service, logger *slog.Logger, corsOrigins []string) http.Handler {
	api := New(svc)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", api.Health)
	mux.HandleFunc("GET /dashboard", api.Dashboard)
	mux.HandleFunc("GET /teams", api.ListTeams)
	mux.HandleFunc("GET /teams/{id}", api.GetTeam)
	mux.HandleFunc("GET /teams/{id}/players", api.TeamPlayers)
	mux.HandleFunc("GET /players/{id}", api.GetPlayer)
	mux.HandleFunc("GET /matches", api.ListMatches)
	mux.HandleFunc("GET /matches/upcoming", api.UpcomingMatches)
	mux.HandleFunc("GET /matches/results", api.MatchResults)
	mux.HandleFunc("GET /standings", api.Standings)
	mux.HandleFunc("GET /stats/topscorers", api.TopScorers)

	var h http.Handler = mux
	h = middleware.StripTrailingSlash(h)
	h = middleware.Timeout(30 * time.Second)(h)
	h = middleware.CORS(corsOrigins)(h)
	h = middleware.Logging(logger)(h)
	h = middleware.Recovery(logger)(h)
	return h
}
