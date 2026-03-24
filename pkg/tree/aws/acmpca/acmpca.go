package acmpca

type PCACertificateAuthority struct {
	CertificateAuthorities []CertificateAuthority `tree:"certificate_authorities"`
}
