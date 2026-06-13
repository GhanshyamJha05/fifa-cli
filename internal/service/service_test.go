package service_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/GhanshyamJha05/fifa-cli/internal/config"
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
	svc, err := service.New(cfg, slog.New(slog.NewTextHandler(os.Stderr, nil)))
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	t.Cleanup(func() { _ = svc.Close() })
	return svc
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

func TestGetStandings(t *testing.T) {
	svc := testService(t)
	standings, err := svc.GetStandings(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(standings) == 0 {
		t.Fatal("expected standings")
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
