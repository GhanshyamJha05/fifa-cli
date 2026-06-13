package provider

import (
	"log/slog"

	"github.com/GhanshyamJha05/fifa-cli/internal/api/apifootball"
	"github.com/GhanshyamJha05/fifa-cli/internal/api/mock"
	"github.com/GhanshyamJha05/fifa-cli/internal/config"
)

// NewFootballProvider builds the configured data provider stack (live/mock + cache).
func NewFootballProvider(cfg *config.Config, store Cacher, logger *slog.Logger) FootballProvider {
	var inner FootballProvider
	if cfg.UseMock || cfg.APIKey == "" {
		logger.Info("using mock data provider (set FIFA_API_KEY for live data)")
		inner = mock.New()
	} else {
		logger.Info("using API-Football provider")
		inner = apifootball.New(cfg.APIBaseURL, cfg.APIKey, cfg.LeagueID, cfg.Season, logger)
	}
	return NewCachedProvider(inner, store, cfg.CacheTTL)
}
