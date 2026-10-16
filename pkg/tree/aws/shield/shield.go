package shield

type Shield struct {
	Subscriptions []Subscription `tree:"subscriptions"`
	Protections   []Protection   `tree:"protections"`
}
