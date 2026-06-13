package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/GhanshyamJha05/fifa-cli/internal/api/mock"
	"github.com/GhanshyamJha05/fifa-cli/internal/config"
	"github.com/GhanshyamJha05/fifa-cli/internal/handlers"
	"github.com/GhanshyamJha05/fifa-cli/internal/provider"
	"github.com/GhanshyamJha05/fifa-cli/internal/repository"
	"github.com/GhanshyamJha05/fifa-cli/internal/service"
	"log/slog"
)

func testRouter(t *testing.T) http.Handler {
	t.Helper()
	cfg := &config.Config{UseMock: true, CacheDir: t.TempDir(), CORSOrigins: []string{"*"}}
	cacheRepo, err := repository.OpenCache(cfg.CacheDir+"/fifa.db", 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cacheRepo.Close() })
	p := provider.NewCachedProvider(mock.New(), cacheRepo, 0)
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	svc := service.NewWithDeps(p, cacheRepo, cfg, logger)
	return handlers.NewRouter(svc, logger, cfg.CORSOrigins)
}

func TestHealth(t *testing.T) {
	h := testRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d", rec.Code)
	}
}

func TestListTeams(t *testing.T) {
	h := testRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/teams?page=1&page_size=5", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	var body struct {
		Items []struct {
			Name string `json:"name"`
		} `json:"items"`
		Total int `json:"total"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Total == 0 {
		t.Fatal("expected teams")
	}
}

func TestGetTeamNotFound(t *testing.T) {
	h := testRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/teams/99999", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status=%d", rec.Code)
	}
}

func TestDashboard(t *testing.T) {
	h := testRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}
