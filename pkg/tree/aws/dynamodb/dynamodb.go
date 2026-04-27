package dynamodb

type DynamoDB struct {
	Tables []Table `tree:"tables"`
}
