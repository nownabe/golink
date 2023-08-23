package backend_test

import (
	"testing"

	"github.com/nownabe/golink/backend"
)

func TestUpdateRedirectCount(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		golink      *backend.Golink
		daysDelayed int
		want        *backend.Golink
	}{
		"no delay": {
			golink: &backend.Golink{
				RedirectCount28Days: 0,
				RedirectCount7Days:  0,
				DailyRedirectCounts: [28]int32{},
			},
			daysDelayed: 0,
			want: &backend.Golink{
				RedirectCount28Days: 1,
				RedirectCount7Days:  1,
				DailyRedirectCounts: [28]int32{0: 1},
			},
		},
		"1 day delay": {
			golink: &backend.Golink{
				RedirectCount28Days: 15,
				RedirectCount7Days:  8,
				DailyRedirectCounts: [28]int32{0: 3, 6: 5, 27: 7},
			},
			daysDelayed: 1,
			want: &backend.Golink{
				RedirectCount28Days: 9,
				RedirectCount7Days:  4,
				DailyRedirectCounts: [28]int32{0: 1, 1: 3, 7: 5},
			},
		},
		"7 days delay": {
			golink: &backend.Golink{
				RedirectCount28Days: 39,
				RedirectCount7Days:  8,
				DailyRedirectCounts: [28]int32{0: 3, 6: 5, 20: 7, 21: 11, 27: 13},
			},
			daysDelayed: 7,
			want: &backend.Golink{
				RedirectCount28Days: 16,
				RedirectCount7Days:  1,
				DailyRedirectCounts: [28]int32{0: 1, 7: 3, 13: 5, 27: 7},
			},
		},
		"28 days delay": {
			golink: &backend.Golink{
				RedirectCount28Days: 3,
				RedirectCount7Days:  3,
				DailyRedirectCounts: [28]int32{0: 1, 1: 2},
			},
			daysDelayed: 28,
			want: &backend.Golink{
				RedirectCount28Days: 1,
				RedirectCount7Days:  1,
				DailyRedirectCounts: [28]int32{0: 1},
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			backend.UpdateRedirectCount(tt.golink, tt.daysDelayed)

			if tt.golink.RedirectCount28Days != tt.want.RedirectCount28Days {
				t.Errorf("RedirectCount28Days: got %d, want %d", tt.golink.RedirectCount28Days, tt.want.RedirectCount28Days)
			}
			if tt.golink.RedirectCount7Days != tt.want.RedirectCount7Days {
				t.Errorf("RedirectCount7Days: got %d, want %d", tt.golink.RedirectCount7Days, tt.want.RedirectCount7Days)
			}
			if tt.golink.DailyRedirectCounts != tt.want.DailyRedirectCounts {
				t.Errorf("DailyRedirectCounts: got %v, want %v", tt.golink.DailyRedirectCounts, tt.want.DailyRedirectCounts)
			}
		})
	}
}
