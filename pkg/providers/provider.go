package providers

import "github.com/infracost/proto/gen/go/infracost/provider"

func ToProto(raw string) provider.Provider {
	switch raw {
	case "aws":
		return provider.Provider_PROVIDER_AWS
	case "google":
		return provider.Provider_PROVIDER_GOOGLE
	case "azurerm":
		return provider.Provider_PROVIDER_AZURERM
	case "kubernetes":
		return provider.Provider_PROVIDER_KUBERNETES
	default:
		return provider.Provider_PROVIDER_UNSPECIFIED
	}
}

func FromProto(raw provider.Provider) string {
	switch raw {
	case provider.Provider_PROVIDER_AWS:
		return "aws"
	case provider.Provider_PROVIDER_GOOGLE:
		return "google"
	case provider.Provider_PROVIDER_AZURERM:
		return "azurerm"
	case provider.Provider_PROVIDER_KUBERNETES:
		return "kubernetes"
	default:
		return ""
	}
}
