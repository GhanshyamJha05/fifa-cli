package provider

import (
	"context"

	"github.com/GhanshyamJha05/fifa-cli/internal/domain"
)

// FootballProvider fetches World Cup data from an external or bundled source.
type FootballProvider interface {
	GetTeams(ctx context.Context) ([]domain.Team, error)
	GetTeam(ctx context.Context, name string) (*domain.Team, error)
	GetTeamByID(ctx context.Context, id int) (*domain.Team, error)
	GetSquad(ctx context.Context, teamName string) ([]domain.Player, error)
	GetPlayer(ctx context.Context, name string) (*domain.Player, error)
	GetPlayerByID(ctx context.Context, id int) (*domain.Player, error)
	GetMatches(ctx context.Context) ([]domain.Match, error)
	GetMatchesToday(ctx context.Context) ([]domain.Match, error)
	GetUpcoming(ctx context.Context, limit int) ([]domain.Match, error)
	GetResults(ctx context.Context) ([]domain.Match, error)
	GetStandings(ctx context.Context) ([]domain.GroupStanding, error)
	GetBracket(ctx context.Context) (*domain.Bracket, error)
	GetStats(ctx context.Context) (*domain.TournamentStats, error)
	GetTournamentInfo(ctx context.Context) (*domain.TournamentInfo, error)
	Search(ctx context.Context, query string) ([]domain.SearchResult, error)
	GetHeadToHead(ctx context.Context, teamA, teamB string) (*domain.HeadToHead, error)
	GetTeamForm(ctx context.Context, teamName string) (*domain.TeamForm, error)
}

// Cache keys for provider responses.
const (
	KeyTeams      = "teams"
	KeyMatches    = "matches"
	KeyStandings  = "standings"
	KeyBracket    = "bracket"
	KeyStats      = "stats"
	KeyTournament = "tournament"
)

// Cacher is the cache interface used by CachedProvider.
type Cacher interface {
	Get(key string, dest any) (bool, error)
	GetStale(key string, dest any) (bool, error)
	Set(key string, data any) error
}
