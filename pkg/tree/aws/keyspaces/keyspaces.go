package keyspaces

type Keyspaces struct {
	Tables []Table `tree:"tables"`
}
