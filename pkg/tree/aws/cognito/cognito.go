package cognito

type Cognito struct {
	UserPools []UserPool `tree:"user_pools"`
}
