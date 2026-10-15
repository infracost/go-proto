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
