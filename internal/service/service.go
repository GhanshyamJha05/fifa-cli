package service

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/GhanshyamJha05/fifa-cli/internal/api"
	"github.com/GhanshyamJha05/fifa-cli/internal/api/apifootball"
	"github.com/GhanshyamJha05/fifa-cli/internal/api/mock"
	"github.com/GhanshyamJha05/fifa-cli/internal/cache"
	"github.com/GhanshyamJha05/fifa-cli/internal/config"
	"github.com/GhanshyamJha05/fifa-cli/internal/domain"
)

// Service is the main application service.
type Service struct {
	provider api.Provider
	cache    *cache.Store
	cfg      *config.Config
	logger   *slog.Logger
}

// New creates a Service from configuration.
func New(cfg *config.Config, logger *slog.Logger) (*Service, error) {
	store, err := cache.Open(filepath.Join(cfg.CacheDir, "fifa.db"), cfg.CacheTTL)
	if err != nil {
		return nil, fmt.Errorf("open cache: %w", err)
	}

	var inner api.Provider
	if cfg.UseMock || cfg.APIKey == "" {
		logger.Info("using mock data provider (set FIFA_API_KEY for live data)")
		inner = mock.New()
	} else {
		logger.Info("using API-Football provider")
		inner = apifootball.New(cfg.APIBaseURL, cfg.APIKey, cfg.LeagueID, cfg.Season, logger)
	}

	provider := api.NewCachedProvider(inner, store, cfg.CacheTTL)

	return &Service{
		provider: provider,
		cache:    store,
		cfg:      cfg,
		logger:   logger,
	}, nil
}

// Close releases resources.
func (s *Service) Close() error {
	return s.cache.Close()
}

// Config returns the service configuration.
func (s *Service) Config() *config.Config { return s.cfg }

func (s *Service) GetTeams(ctx context.Context) ([]domain.Team, error) {
	return s.provider.GetTeams(ctx)
}

func (s *Service) GetTeam(ctx context.Context, name string) (*domain.Team, error) {
	return s.provider.GetTeam(ctx, name)
}

func (s *Service) GetSquad(ctx context.Context, teamName string) ([]domain.Player, error) {
	return s.provider.GetSquad(ctx, teamName)
}

func (s *Service) GetPlayer(ctx context.Context, name string) (*domain.Player, error) {
	return s.provider.GetPlayer(ctx, name)
}

func (s *Service) GetMatches(ctx context.Context) ([]domain.Match, error) {
	return s.provider.GetMatches(ctx)
}

func (s *Service) GetMatchesToday(ctx context.Context) ([]domain.Match, error) {
	return s.provider.GetMatchesToday(ctx)
}

func (s *Service) GetUpcoming(ctx context.Context, limit int) ([]domain.Match, error) {
	return s.provider.GetUpcoming(ctx, limit)
}

func (s *Service) GetResults(ctx context.Context) ([]domain.Match, error) {
	return s.provider.GetResults(ctx)
}

func (s *Service) GetStandings(ctx context.Context) ([]domain.GroupStanding, error) {
	return s.provider.GetStandings(ctx)
}

func (s *Service) GetBracket(ctx context.Context) (*domain.Bracket, error) {
	return s.provider.GetBracket(ctx)
}

func (s *Service) GetStats(ctx context.Context) (*domain.TournamentStats, error) {
	return s.provider.GetStats(ctx)
}

func (s *Service) GetTournamentInfo(ctx context.Context) (*domain.TournamentInfo, error) {
	return s.provider.GetTournamentInfo(ctx)
}

func (s *Service) Search(ctx context.Context, query string) ([]domain.SearchResult, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}
	return s.provider.Search(ctx, query)
}

func (s *Service) GetHeadToHead(ctx context.Context, teamA, teamB string) (*domain.HeadToHead, error) {
	return s.provider.GetHeadToHead(ctx, teamA, teamB)
}

func (s *Service) GetTeamForm(ctx context.Context, teamName string) (*domain.TeamForm, error) {
	return s.provider.GetTeamForm(ctx, teamName)
}

// FilterMatchesByDate returns matches on a given date string (YYYY-MM-DD).
func FilterMatchesByDate(matches []domain.Match, dateStr string) []domain.Match {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil
	}
	var out []domain.Match
	for _, m := range matches {
		if m.Date.Year() == t.Year() && m.Date.Month() == t.Month() && m.Date.Day() == t.Day() {
			out = append(out, m)
		}
	}
	return out
}

// TournamentProgress returns completion percentage.
func TournamentProgress(info *domain.TournamentInfo) float64 {
	if info.TotalMatches == 0 {
		return 0
	}
	return float64(info.Completed) / float64(info.TotalMatches) * 100
}
