package sns

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Subscription struct {
	resource.Resource `tree:"-"`
	TopicARN          value.String                       `tree:"topic_arn"`
	Protocol          value.Value[SubscriptionProtocol] `tree:"protocol"`
}

type SubscriptionProtocol uint32

const (
	SubscriptionProtocolUnknown     SubscriptionProtocol = iota
	SubscriptionProtocolHTTP
	SubscriptionProtocolHTTPS
	SubscriptionProtocolEmail
	SubscriptionProtocolEmailJSON
	SubscriptionProtocolSMS
	SubscriptionProtocolSQS
	SubscriptionProtocolApplication
	SubscriptionProtocolLambda
	SubscriptionProtocolFirehose
)
