package service

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/GhanshyamJha05/fifa-cli/internal/config"
	"github.com/GhanshyamJha05/fifa-cli/internal/domain"
	"github.com/GhanshyamJha05/fifa-cli/internal/provider"
	"github.com/GhanshyamJha05/fifa-cli/internal/repository"
	"golang.org/x/sync/errgroup"
)

// Service is the main application service with injected dependencies.
type Service struct {
	provider provider.FootballProvider
	cache    *repository.CacheRepository
	cfg      *config.Config
	logger   *slog.Logger

	refreshMu sync.Mutex
	refreshing bool
}

// New creates a Service from configuration (composition root).
func New(cfg *config.Config, logger *slog.Logger) (*Service, error) {
	cacheRepo, err := repository.OpenCache(filepath.Join(cfg.CacheDir, "fifa.db"), cfg.CacheTTL)
	if err != nil {
		return nil, fmt.Errorf("open cache: %w", err)
	}
	p := provider.NewFootballProvider(cfg, cacheRepo, logger)
	return &Service{
		provider: p,
		cache:    cacheRepo,
		cfg:      cfg,
		logger:   logger,
	}, nil
}

// NewWithDeps injects dependencies for testing.
func NewWithDeps(p provider.FootballProvider, cache *repository.CacheRepository, cfg *config.Config, logger *slog.Logger) *Service {
	return &Service{provider: p, cache: cache, cfg: cfg, logger: logger}
}

// Close releases resources.
func (s *Service) Close() error {
	if s.cache == nil {
		return nil
	}
	return s.cache.Close()
}

// Config returns the service configuration.
func (s *Service) Config() *config.Config { return s.cfg }

// Provider returns the underlying football provider (for advanced use).
func (s *Service) Provider() provider.FootballProvider { return s.provider }

// DashboardData holds concurrently loaded home screen data.
type DashboardData struct {
	Info          *domain.TournamentInfo
	TodayMatches  []domain.Match
	Teams         []domain.Team
	Standings     []domain.GroupStanding
	Stats         *domain.TournamentStats
}

// LoadDashboard fetches tournament overview data concurrently.
func (s *Service) LoadDashboard(ctx context.Context) (*DashboardData, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	var data DashboardData
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		info, err := s.provider.GetTournamentInfo(ctx)
		if err != nil {
			return err
		}
		data.Info = info
		return nil
	})
	g.Go(func() error {
		matches, err := s.provider.GetMatchesToday(ctx)
		if err != nil {
			return err
		}
		data.TodayMatches = matches
		return nil
	})
	g.Go(func() error {
		teams, err := s.provider.GetTeams(ctx)
		if err != nil {
			return err
		}
		data.Teams = teams
		return nil
	})
	g.Go(func() error {
		standings, err := s.provider.GetStandings(ctx)
		if err != nil {
			return err
		}
		data.Standings = standings
		return nil
	})
	g.Go(func() error {
		stats, err := s.provider.GetStats(ctx)
		if err != nil {
			return err
		}
		data.Stats = stats
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return &data, nil
}

// RefreshCache warms core caches in the background without blocking the caller.
func (s *Service) RefreshCache(parent context.Context) {
	s.refreshMu.Lock()
	if s.refreshing {
		s.refreshMu.Unlock()
		return
	}
	s.refreshing = true
	s.refreshMu.Unlock()

	go func() {
		defer func() {
			s.refreshMu.Lock()
			s.refreshing = false
			s.refreshMu.Unlock()
		}()

		ctx, cancel := context.WithTimeout(parent, 30*time.Second)
		defer cancel()

		g, ctx := errgroup.WithContext(ctx)
		fetches := []func() error{
			func() error { _, err := s.provider.GetTeams(ctx); return err },
			func() error { _, err := s.provider.GetMatches(ctx); return err },
			func() error { _, err := s.provider.GetStandings(ctx); return err },
			func() error { _, err := s.provider.GetStats(ctx); return err },
		}
		for _, fn := range fetches {
			g.Go(fn)
		}
		if err := g.Wait(); err != nil {
			s.logger.Warn("background cache refresh failed", "error", err)
		} else {
			s.logger.Debug("background cache refresh completed")
		}
	}()
}

func (s *Service) GetTeams(ctx context.Context) ([]domain.Team, error) {
	return s.provider.GetTeams(ctx)
}

func (s *Service) GetTeam(ctx context.Context, name string) (*domain.Team, error) {
	return s.provider.GetTeam(ctx, name)
}

func (s *Service) GetTeamByID(ctx context.Context, id int) (*domain.Team, error) {
	return s.provider.GetTeamByID(ctx, id)
}

func (s *Service) GetSquad(ctx context.Context, teamName string) ([]domain.Player, error) {
	return s.provider.GetSquad(ctx, teamName)
}

func (s *Service) GetPlayer(ctx context.Context, name string) (*domain.Player, error) {
	return s.provider.GetPlayer(ctx, name)
}

func (s *Service) GetPlayerByID(ctx context.Context, id int) (*domain.Player, error) {
	return s.provider.GetPlayerByID(ctx, id)
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

// SearchPlayers filters players by name (case-insensitive) from cached teams/squads path.
func SearchPlayers(players []domain.Player, query string) []domain.Player {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return players
	}
	var out []domain.Player
	for _, p := range players {
		if strings.Contains(strings.ToLower(p.Name), q) {
			out = append(out, p)
		}
	}
	return out
}
