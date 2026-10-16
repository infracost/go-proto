package qbusiness

type QBusiness struct {
	Applications []Application `tree:"applications"`
	Indexes      []Index       `tree:"indexes"`
}
