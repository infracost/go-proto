package quicksight

type QuickSight struct {
	Users                []User                `tree:"users"`
	AccountSubscriptions []AccountSubscription `tree:"account_subscriptions"`
}
