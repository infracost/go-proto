package apigateway

type APIGateway struct {
	RestAPIs []RestAPI `tree:"rest_apis"`
	Stages   []Stage   `tree:"stages"`
}
