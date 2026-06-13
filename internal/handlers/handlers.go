package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/GhanshyamJha05/fifa-cli/internal/service"
	"github.com/GhanshyamJha05/fifa-cli/pkg/httputil"
	"github.com/GhanshyamJha05/fifa-cli/pkg/pagination"
)

// API exposes HTTP handlers for the FIFA service.
type API struct {
	svc *service.Service
}

// New creates HTTP handlers.
func New(svc *service.Service) *API {
	return &API{svc: svc}
}

func (a *API) ctx(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), 15*time.Second)
}

func queryInt(r *http.Request, key string, defaultVal int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 {
		return defaultVal
	}
	return n
}

func pathID(r *http.Request, key string) (int, error) {
	v := r.PathValue(key)
	return strconv.Atoi(v)
}

// Health handles GET /health
func (a *API) Health(w http.ResponseWriter, r *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

// ListTeams handles GET /teams
func (a *API) ListTeams(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := a.ctx(r)
	defer cancel()

	teams, err := a.svc.GetTeams(ctx)
	if err != nil {
		httputil.Error(w, http.StatusBadGateway, err.Error())
		return
	}
	page := queryInt(r, "page", 1)
	size := queryInt(r, "page_size", 50)
	httputil.JSON(w, http.StatusOK, pagination.Paginate(teams, page, size))
}

// GetTeam handles GET /teams/{id}
func (a *API) GetTeam(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := a.ctx(r)
	defer cancel()

	id, err := pathID(r, "id")
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid team id")
		return
	}
	team, err := a.svc.GetTeamByID(ctx, id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, team)
}

// TeamPlayers handles GET /teams/{id}/players
func (a *API) TeamPlayers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := a.ctx(r)
	defer cancel()

	id, err := pathID(r, "id")
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid team id")
		return
	}
	team, err := a.svc.GetTeamByID(ctx, id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, err.Error())
		return
	}
	players, err := a.svc.GetSquad(ctx, team.Name)
	if err != nil {
		httputil.Error(w, http.StatusBadGateway, err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, players)
}

// GetPlayer handles GET /players/{id}
func (a *API) GetPlayer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := a.ctx(r)
	defer cancel()

	id, err := pathID(r, "id")
	if err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid player id")
		return
	}
	player, err := a.svc.GetPlayerByID(ctx, id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, player)
}

// ListMatches handles GET /matches
func (a *API) ListMatches(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := a.ctx(r)
	defer cancel()

	matches, err := a.svc.GetMatches(ctx)
	if err != nil {
		httputil.Error(w, http.StatusBadGateway, err.Error())
		return
	}
	if date := r.URL.Query().Get("date"); date != "" {
		matches = service.FilterMatchesByDate(matches, date)
	}
	page := queryInt(r, "page", 1)
	size := queryInt(r, "page_size", 50)
	httputil.JSON(w, http.StatusOK, pagination.Paginate(matches, page, size))
}

// UpcomingMatches handles GET /matches/upcoming
func (a *API) UpcomingMatches(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := a.ctx(r)
	defer cancel()

	limit := queryInt(r, "limit", 20)
	matches, err := a.svc.GetUpcoming(ctx, limit)
	if err != nil {
		httputil.Error(w, http.StatusBadGateway, err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, matches)
}

// MatchResults handles GET /matches/results
func (a *API) MatchResults(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := a.ctx(r)
	defer cancel()

	matches, err := a.svc.GetResults(ctx)
	if err != nil {
		httputil.Error(w, http.StatusBadGateway, err.Error())
		return
	}
	page := queryInt(r, "page", 1)
	size := queryInt(r, "page_size", 50)
	httputil.JSON(w, http.StatusOK, pagination.Paginate(matches, page, size))
}

// Standings handles GET /standings
func (a *API) Standings(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := a.ctx(r)
	defer cancel()

	standings, err := a.svc.GetStandings(ctx)
	if err != nil {
		httputil.Error(w, http.StatusBadGateway, err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, standings)
}

// TopScorers handles GET /stats/topscorers
func (a *API) TopScorers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := a.ctx(r)
	defer cancel()

	stats, err := a.svc.GetStats(ctx)
	if err != nil {
		httputil.Error(w, http.StatusBadGateway, err.Error())
		return
	}
	limit := queryInt(r, "limit", 20)
	scorers := stats.TopScorers
	if limit < len(scorers) {
		scorers = scorers[:limit]
	}
	httputil.JSON(w, http.StatusOK, scorers)
}

// Dashboard handles GET /dashboard (concurrent load).
func (a *API) Dashboard(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := a.ctx(r)
	defer cancel()

	data, err := a.svc.LoadDashboard(ctx)
	if err != nil {
		httputil.Error(w, http.StatusBadGateway, err.Error())
		return
	}
	httputil.JSON(w, http.StatusOK, data)
}
