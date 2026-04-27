package cloudhsmv2

type CloudHSMV2 struct {
	HSMs []HSM `tree:"hsms"`
}
