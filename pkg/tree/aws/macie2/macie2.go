package macie2

type Macie2 struct {
	Accounts           []Account           `tree:"accounts"`
	ClassificationJobs []ClassificationJob `tree:"classification_jobs"`
}
