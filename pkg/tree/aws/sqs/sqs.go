package sqs

type SQS struct {
	Queues []Queue `tree:"queues"`
}
