package identityplatform

type IdentityPlatform struct {
	FederatedIdentityCredentials []FederatedIdentityCredential `tree:"federated_identity_credentials"`
}
