package composer

type Composer struct {
	Environments []Environment `tree:"environments"`
}
