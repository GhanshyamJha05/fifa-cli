package cache_test

import (
	"testing"
	"time"

	"github.com/GhanshyamJha05/fifa-cli/internal/cache"
)

func TestStoreSetGet(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	store, err := cache.Open(dir+"/test.db", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	type payload struct {
		Name string `json:"name"`
	}
	in := payload{Name: "Brazil"}
	if err := store.Set("team", in); err != nil {
		t.Fatal(err)
	}

	var out payload
	ok, err := store.Get("team", &out)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected cache hit")
	}
	if out.Name != "Brazil" {
		t.Fatalf("got %q", out.Name)
	}
}

func TestStoreExpired(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	store, err := cache.Open(dir+"/test.db", time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	if err := store.Set("k", "v"); err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Millisecond)

	var out string
	ok, err := store.Get("k", &out)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected cache miss after TTL")
	}
}

func TestStoreStale(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	store, err := cache.Open(dir+"/test.db", time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	if err := store.Set("k", "stale"); err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Millisecond)

	var out string
	ok, err := store.GetStale("k", &out)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || out != "stale" {
		t.Fatalf("stale get failed: ok=%v out=%q", ok, out)
	}
}

func BenchmarkCacheSetGet(b *testing.B) {
	dir := b.TempDir()
	store, err := cache.Open(dir+"/bench.db", time.Hour)
	if err != nil {
		b.Fatal(err)
	}
	defer store.Close()

	data := map[string]int{"goals": 42}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = store.Set("bench", data)
		var out map[string]int
		_, _ = store.Get("bench", &out)
	}
}
