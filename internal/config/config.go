package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Config holds application configuration.
type Config struct {
	APIKey       string        `mapstructure:"api_key"`
	APIBaseURL   string        `mapstructure:"api_base_url"`
	LeagueID     int           `mapstructure:"league_id"`
	Season       int           `mapstructure:"season"`
	CacheDir     string        `mapstructure:"cache_dir"`
	CacheTTL     time.Duration `mapstructure:"cache_ttl"`
	UseMock      bool          `mapstructure:"use_mock"`
	Theme        string        `mapstructure:"theme"`
	Favorites    Favorites     `mapstructure:"favorites"`
	LogLevel     string        `mapstructure:"log_level"`
}

// Favorites stores user favorite teams and players.
type Favorites struct {
	Teams   []string `mapstructure:"teams"`
	Players []string `mapstructure:"players"`
}

// Load reads configuration from env, config file, and defaults.
func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("home dir: %w", err)
	}

	configDir := filepath.Join(home, ".fifa-cli")
	v.AddConfigPath(configDir)
	v.AddConfigPath(".")

	v.SetEnvPrefix("FIFA")
	v.AutomaticEnv()
	v.BindEnv("api_key")
	v.BindEnv("use_mock")

	v.SetDefault("api_base_url", "https://v3.football.api-sports.io")
	v.SetDefault("league_id", 1)
	v.SetDefault("season", 2026)
	v.SetDefault("cache_dir", filepath.Join(configDir, "cache"))
	v.SetDefault("cache_ttl", 15*time.Minute)
	v.SetDefault("use_mock", true)
	v.SetDefault("theme", "dark")
	v.SetDefault("log_level", "info")

	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return nil, fmt.Errorf("create config dir: %w", err)
	}

	_ = v.ReadInConfig()

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if cfg.APIKey == "" {
		cfg.UseMock = true
	}

	if err := os.MkdirAll(cfg.CacheDir, 0o755); err != nil {
		return nil, fmt.Errorf("create cache dir: %w", err)
	}

	return cfg, nil
}
