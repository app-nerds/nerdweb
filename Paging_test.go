package nerdweb_test

import (
	"testing"

	"github.com/app-nerds/nerdweb/v2"
)

func TestAdjustPage(t *testing.T) {
	tests := []struct {
		name string
		want int
		page int
	}{
		{
			name: "Returns the value of page minus one",
			want: 5,
			page: 6,
		},
		{
			name: "Returns zero when page is less than zero",
			want: 0,
			page: -2,
		},
	}

	for _, tt := range tests {
		got := nerdweb.AdjustPage(tt.page)

		if got != tt.want {
			t.Errorf("wanted %d, got %d", tt.want, got)
		}
	}
}

func TestHasNextPage(t *testing.T) {
	type args struct {
		page        int
		pageSize    int
		recordCount int
	}

	tests := []struct {
		name string
		want bool
		args args
	}{
		{
			name: "Returns true when there are more pages",
			want: true,
			args: args{
				page:        1,
				pageSize:    50,
				recordCount: 150,
			},
		},
		{
			name: "Returns false when there are no more pages",
			want: false,
			args: args{
				page:        3,
				pageSize:    50,
				recordCount: 150,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := nerdweb.HasNextPage(tt.args.page, tt.args.pageSize, tt.args.recordCount)

			if got != tt.want {
				t.Errorf("want %v, got %v", tt.want, got)
			}
		})
	}
}

func TestTotalPages(t *testing.T) {
	type args struct {
		pageSize    int
		recordCount int
	}

	tests := []struct {
		name string
		want int
		args args
	}{
		{
			name: "Returns correct value with even record count",
			want: 2,
			args: args{
				pageSize:    50,
				recordCount: 100,
			},
		},
		{
			name: "Returns correct value with odd record count",
			want: 3,
			args: args{
				pageSize:    50,
				recordCount: 101,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := nerdweb.TotalPages(tt.args.pageSize, tt.args.recordCount)

			if got != tt.want {
				t.Errorf("want %d, got %d", tt.want, got)
			}
		})
	}
}
