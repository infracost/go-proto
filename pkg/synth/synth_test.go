package synth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPattern(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"prefix", "cloudtrail-" + Marker(KindPrefix, "ab12.cd34"), "cloudtrail-*"},
		{"suffix", Identifier + "-api", "*-api"},
		{"middle", "coast-" + Marker("random_pet", "9f3a") + "-web", "coast-*-web"},
		{"upper mock", "INFRACOST-MOCK-44956BE29F34-api", "*-api"},
		{"upper marker", "DATA-%IC:RANDOM_PET:9F3A%", "DATA-*"},
		{"with length", "data-" + MarkerLen("random_string", "ab.cd", 6), "data-*"},
		{"adjacent", Identifier + Marker("random_id", "1") + "-api", "*-api"},
		{"too short", "a-" + Marker("random_pet", "1"), ""},
		{"only wildcard", Identifier, ""},
		{"no unknown part", "my-bucket", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Pattern(tt.in))
		})
	}
}

func TestMarkers(t *testing.T) {
	s := "a-" + Marker("random_pet", "x.y") + "-" + MarkerLen("random_string", "z", 6)
	got := Markers(s)
	assert.Len(t, got, 2)
	assert.Equal(t, Marker("random_pet", "x.y"), s[got[0].Start:got[0].End])
	assert.Equal(t, DefaultLength, got[0].Length)
	assert.Equal(t, MarkerLen("random_string", "z", 6), s[got[1].Start:got[1].End])
	assert.Equal(t, 6, got[1].Length)
}
