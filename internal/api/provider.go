package api

import (
	"context"
	"time"

	"github.com/GhanshyamJha05/fifa-cli/internal/domain"
)

// Provider fetches World Cup data from an external source.
type Provider interface {
	GetTeams(ctx context.Context) ([]domain.Team, error)
	GetTeam(ctx context.Context, name string) (*domain.Team, error)
	GetSquad(ctx context.Context, teamName string) ([]domain.Player, error)
	GetPlayer(ctx context.Context, name string) (*domain.Player, error)
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

// CachedProvider wraps a Provider with local caching.
type CachedProvider struct {
	inner Provider
	store Cacher
	ttl   time.Duration
}

// Cacher is the cache interface used by CachedProvider.
type Cacher interface {
	Get(key string, dest any) (bool, error)
	GetStale(key string, dest any) (bool, error)
	Set(key string, data any) error
}

// NewCachedProvider wraps a provider with caching.
func NewCachedProvider(inner Provider, store Cacher, ttl time.Duration) *CachedProvider {
	return &CachedProvider{inner: inner, store: store, ttl: ttl}
}

func (c *CachedProvider) GetTeams(ctx context.Context) ([]domain.Team, error) {
	var teams []domain.Team
	if ok, _ := c.store.Get(KeyTeams, &teams); ok {
		return teams, nil
	}
	teams, err := c.inner.GetTeams(ctx)
	if err != nil {
		if ok, _ := c.store.GetStale(KeyTeams, &teams); ok {
			return teams, nil
		}
		return nil, err
	}
	_ = c.store.Set(KeyTeams, teams)
	return teams, nil
}

func (c *CachedProvider) GetTeam(ctx context.Context, name string) (*domain.Team, error) {
	return c.inner.GetTeam(ctx, name)
}

func (c *CachedProvider) GetSquad(ctx context.Context, teamName string) ([]domain.Player, error) {
	return c.inner.GetSquad(ctx, teamName)
}

func (c *CachedProvider) GetPlayer(ctx context.Context, name string) (*domain.Player, error) {
	return c.inner.GetPlayer(ctx, name)
}

func (c *CachedProvider) GetMatches(ctx context.Context) ([]domain.Match, error) {
	var matches []domain.Match
	if ok, _ := c.store.Get(KeyMatches, &matches); ok {
		return matches, nil
	}
	matches, err := c.inner.GetMatches(ctx)
	if err != nil {
		if ok, _ := c.store.GetStale(KeyMatches, &matches); ok {
			return matches, nil
		}
		return nil, err
	}
	_ = c.store.Set(KeyMatches, matches)
	return matches, nil
}

func (c *CachedProvider) GetMatchesToday(ctx context.Context) ([]domain.Match, error) {
	return c.inner.GetMatchesToday(ctx)
}

func (c *CachedProvider) GetUpcoming(ctx context.Context, limit int) ([]domain.Match, error) {
	return c.inner.GetUpcoming(ctx, limit)
}

func (c *CachedProvider) GetResults(ctx context.Context) ([]domain.Match, error) {
	return c.inner.GetResults(ctx)
}

func (c *CachedProvider) GetStandings(ctx context.Context) ([]domain.GroupStanding, error) {
	var standings []domain.GroupStanding
	if ok, _ := c.store.Get(KeyStandings, &standings); ok {
		return standings, nil
	}
	standings, err := c.inner.GetStandings(ctx)
	if err != nil {
		if ok, _ := c.store.GetStale(KeyStandings, &standings); ok {
			return standings, nil
		}
		return nil, err
	}
	_ = c.store.Set(KeyStandings, standings)
	return standings, nil
}

func (c *CachedProvider) GetBracket(ctx context.Context) (*domain.Bracket, error) {
	var bracket domain.Bracket
	if ok, _ := c.store.Get(KeyBracket, &bracket); ok {
		return &bracket, nil
	}
	b, err := c.inner.GetBracket(ctx)
	if err != nil {
		if ok, _ := c.store.GetStale(KeyBracket, &bracket); ok {
			return &bracket, nil
		}
		return nil, err
	}
	_ = c.store.Set(KeyBracket, b)
	return b, nil
}

func (c *CachedProvider) GetStats(ctx context.Context) (*domain.TournamentStats, error) {
	var stats domain.TournamentStats
	if ok, _ := c.store.Get(KeyStats, &stats); ok {
		return &stats, nil
	}
	s, err := c.inner.GetStats(ctx)
	if err != nil {
		if ok, _ := c.store.GetStale(KeyStats, &stats); ok {
			return &stats, nil
		}
		return nil, err
	}
	_ = c.store.Set(KeyStats, s)
	return s, nil
}

func (c *CachedProvider) GetTournamentInfo(ctx context.Context) (*domain.TournamentInfo, error) {
	var info domain.TournamentInfo
	if ok, _ := c.store.Get(KeyTournament, &info); ok {
		return &info, nil
	}
	i, err := c.inner.GetTournamentInfo(ctx)
	if err != nil {
		if ok, _ := c.store.GetStale(KeyTournament, &info); ok {
			return &info, nil
		}
		return nil, err
	}
	_ = c.store.Set(KeyTournament, i)
	return i, nil
}

func (c *CachedProvider) Search(ctx context.Context, query string) ([]domain.SearchResult, error) {
	return c.inner.Search(ctx, query)
}

func (c *CachedProvider) GetHeadToHead(ctx context.Context, teamA, teamB string) (*domain.HeadToHead, error) {
	return c.inner.GetHeadToHead(ctx, teamA, teamB)
}

func (c *CachedProvider) GetTeamForm(ctx context.Context, teamName string) (*domain.TeamForm, error) {
	return c.inner.GetTeamForm(ctx, teamName)
}
