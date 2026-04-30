package appservice

type SSLState uint32

const (
	SSLStateUnknown SSLState = iota
	SSLStateIPBasedEnabled
	SSLStateSniEnabled
)

type CertificateProductType uint32

const (
	CertificateProductTypeUnknown CertificateProductType = iota
	CertificateProductTypeStandard
	CertificateProductTypeWildCard
)
