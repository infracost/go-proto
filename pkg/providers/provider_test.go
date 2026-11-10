package providers

import (
	"testing"

	"github.com/infracost/proto/gen/go/infracost/provider"
	"github.com/stretchr/testify/assert"
)

func TestRoundTrip(t *testing.T) {
	for raw, want := range map[string]provider.Provider{
		"aws":        provider.Provider_PROVIDER_AWS,
		"google":     provider.Provider_PROVIDER_GOOGLE,
		"azurerm":    provider.Provider_PROVIDER_AZURERM,
		"kubernetes": provider.Provider_PROVIDER_KUBERNETES,
		"ai":         provider.Provider_PROVIDER_AI,
	} {
		t.Run(raw, func(t *testing.T) {
			assert.Equal(t, want, ToProto(raw))
			assert.Equal(t, raw, FromProto(want))
		})
	}
}

func TestUnknown(t *testing.T) {
	assert.Equal(t, provider.Provider_PROVIDER_UNSPECIFIED, ToProto("unknown"))
	assert.Empty(t, FromProto(provider.Provider_PROVIDER_UNSPECIFIED))
}
