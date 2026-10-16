package appsync

type AppSync struct {
	APICaches   []APICache   `tree:"api_caches"`
	GraphQLAPIs []GraphQLAPI `tree:"graphql_apis"`
	EventAPIs   []EventAPI   `tree:"event_apis"`
}
