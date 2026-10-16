package bedrock

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Guardrail struct {
	resource.Resource          `tree:"-"`
	ContentPolicy              value.Bool `tree:"content_policy"`
	TopicPolicy                value.Bool `tree:"topic_policy"`
	SensitiveInformationPolicy value.Bool `tree:"sensitive_information_policy"`
	ContextualGroundingPolicy  value.Bool `tree:"contextual_grounding_policy"`
}
