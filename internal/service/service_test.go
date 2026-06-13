package service_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/GhanshyamJha05/fifa-cli/internal/api/mock"
	"github.com/GhanshyamJha05/fifa-cli/internal/config"
	"github.com/GhanshyamJha05/fifa-cli/internal/domain"
	"github.com/GhanshyamJha05/fifa-cli/internal/provider"
	"github.com/GhanshyamJha05/fifa-cli/internal/repository"
	"github.com/GhanshyamJha05/fifa-cli/internal/service"
)

func testService(t *testing.T) *service.Service {
	t.Helper()
	cfg := &config.Config{
		UseMock:  true,
		CacheDir: t.TempDir(),
		CacheTTL: 0,
		Theme:    "dark",
	}
	cacheRepo, err := repository.OpenCache(cfg.CacheDir+"/fifa.db", cfg.CacheTTL)
	if err != nil {
		t.Fatalf("cache: %v", err)
	}
	t.Cleanup(func() { _ = cacheRepo.Close() })

	p := provider.NewCachedProvider(mock.New(), cacheRepo, cfg.CacheTTL)
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	return service.NewWithDeps(p, cacheRepo, cfg, logger)
}

func TestGetTeams(t *testing.T) {
	svc := testService(t)
	teams, err := svc.GetTeams(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(teams) == 0 {
		t.Fatal("expected teams")
	}
}

func TestGetTeam(t *testing.T) {
	svc := testService(t)
	team, err := svc.GetTeam(context.Background(), "Brazil")
	if err != nil {
		t.Fatal(err)
	}
	if team.Name != "Brazil" {
		t.Fatalf("got %q", team.Name)
	}
}

func TestGetTeamByID(t *testing.T) {
	svc := testService(t)
	team, err := svc.GetTeamByID(context.Background(), 5)
	if err != nil {
		t.Fatal(err)
	}
	if team.Name != "Brazil" {
		t.Fatalf("got %q", team.Name)
	}
}

func TestSearch(t *testing.T) {
	svc := testService(t)
	results, err := svc.Search(context.Background(), "Messi")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("expected search results for Messi")
	}
}

func TestLoadDashboardConcurrent(t *testing.T) {
	svc := testService(t)
	data, err := svc.LoadDashboard(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if data.Info == nil {
		t.Fatal("expected tournament info")
	}
	if len(data.Teams) == 0 {
		t.Fatal("expected teams in dashboard")
	}
}

func TestTournamentProgress(t *testing.T) {
	info, err := testService(t).GetTournamentInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	pct := service.TournamentProgress(info)
	if pct < 0 || pct > 100 {
		t.Fatalf("invalid progress: %f", pct)
	}
}

func TestFilterMatchesByDate(t *testing.T) {
	svc := testService(t)
	matches, err := svc.GetMatches(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) == 0 {
		t.Fatal("no matches")
	}
	date := matches[0].Date.Format("2006-01-02")
	filtered := service.FilterMatchesByDate(matches, date)
	if len(filtered) == 0 {
		t.Fatal("expected matches on date")
	}
}

func TestSearchPlayers(t *testing.T) {
	players := []domain.Player{
		{Name: "Lionel Messi"},
		{Name: "Neymar"},
		{Name: "Harry Kane"},
	}
	results := service.SearchPlayers(players, "messi")
	if len(results) != 1 || results[0].Name != "Lionel Messi" {
		t.Fatalf("unexpected results: %+v", results)
	}
}

func BenchmarkSearchPlayers(b *testing.B) {
	players := make([]domain.Player, 100)
	for i := range players {
		players[i] = domain.Player{Name: "Player " + string(rune('A'+i%26))}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.SearchPlayers(players, "Player")
	}
}
