package transfer

type Transfer struct {
	Servers []Server `tree:"servers"`
}
