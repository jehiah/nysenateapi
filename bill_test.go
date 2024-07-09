package nysenateapi

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Test_parseTime(t *testing.T) {
	type testCase struct {
		s    string
		want time.Time
	}
	tests := []testCase{
		{
			s:    "2024-05-06",
			want: time.Date(2024, 5, 6, 0, 0, 0, 0, time.UTC),
		},
		{
			s:    "2024-05-06T00:00",
			want: time.Date(2024, 5, 6, 0, 0, 0, 0, time.UTC),
		},
		{
			s:    "2022-12-28T11:49:26.931773",
			want: time.Date(2022, 12, 28, 11, 49, 26, 931773000, time.UTC),
		},
	}
	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if got := parseTime(tc.s); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("parseTime() = %v, want %v", got, tc.want)
			}
		})
	}
}
