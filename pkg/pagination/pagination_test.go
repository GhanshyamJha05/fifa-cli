package pagination_test

import (
	"testing"

	"github.com/GhanshyamJha05/fifa-cli/pkg/pagination"
)

func TestPaginate(t *testing.T) {
	t.Parallel()
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	tests := []struct {
		name       string
		page       int
		size       int
		wantLen    int
		wantTotal  int
		wantPages  int
	}{
		{"first page", 1, 3, 3, 10, 4},
		{"last page partial", 4, 3, 1, 10, 4},
		{"default page", 0, 0, 10, 10, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := pagination.Paginate(items, tt.page, tt.size)
			if len(r.Items) != tt.wantLen {
				t.Fatalf("items len=%d want %d", len(r.Items), tt.wantLen)
			}
			if r.Total != tt.wantTotal {
				t.Fatalf("total=%d want %d", r.Total, tt.wantTotal)
			}
			if r.TotalPages != tt.wantPages {
				t.Fatalf("pages=%d want %d", r.TotalPages, tt.wantPages)
			}
		})
	}
}

func BenchmarkPaginate(b *testing.B) {
	items := make([]int, 1000)
	for i := range items {
		items[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pagination.Paginate(items, 2, 50)
	}
}
