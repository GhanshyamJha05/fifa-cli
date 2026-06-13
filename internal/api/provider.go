// Package api provides backward-compatible aliases for the provider layer.
package api

import "github.com/GhanshyamJha05/fifa-cli/internal/provider"

// Provider is an alias for FootballProvider.
type Provider = provider.FootballProvider

// Cacher is an alias for the cache interface.
type Cacher = provider.Cacher

// Cache keys re-exported for compatibility.
const (
	KeyTeams      = provider.KeyTeams
	KeyMatches    = provider.KeyMatches
	KeyStandings  = provider.KeyStandings
	KeyBracket    = provider.KeyBracket
	KeyStats      = provider.KeyStats
	KeyTournament = provider.KeyTournament
)

// CachedProvider re-export.
type CachedProvider = provider.CachedProvider

// NewCachedProvider wraps a provider with caching.
var NewCachedProvider = provider.NewCachedProvider
