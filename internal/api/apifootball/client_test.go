package apifootball_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/GhanshyamJha05/fifa-cli/internal/api/apifootball"
)

func TestClientGetTeamsSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path != "/teams" {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"response": []map[string]any{
				{"team": map[string]any{"id": 1, "name": "Brazil", "code": "BRA"}},
			},
		})
	}))
	defer srv.Close()

	client := apifootball.New(srv.URL, "test-key", 1, 2026, slog.New(slog.NewTextHandler(os.Stderr, nil)))
	teams, err := client.GetTeams(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(teams) != 1 || teams[0].Name != "Brazil" {
		t.Fatalf("unexpected teams: %+v", teams)
	}
}

func TestClientServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := apifootball.New(srv.URL, "test-key", 1, 2026, slog.New(slog.NewTextHandler(os.Stderr, nil)))
	_, err := client.GetTeams(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestClientInvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("not-json"))
	}))
	defer srv.Close()

	client := apifootball.New(srv.URL, "test-key", 1, 2026, slog.New(slog.NewTextHandler(os.Stderr, nil)))
	_, err := client.GetTeams(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid json")
	}
}
