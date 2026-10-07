// Package ai is the AI provider in the resource tree. It holds where and how
// AI models are used, as opposed to the cloud resources that host them.
//
// The code service holds LLM call sites found in application source code by
// the appcode parser plugin: the SDK clients, the calls made through them, and
// the model IDs those calls use.
//
// AI is a provider only in the tree's sense: a top-level section alongside aws
// and azure, filled in by a parser plugin like the others. Unlike the cloud
// sections, it has no matching provider plugin, so nothing prices these
// resources or evaluates FinOps policies against them. Provider plugins are
// handed the whole tree, but each reads only its own section and skips this
// one. Consumers read these resources straight from the tree instead.
//
// For the same reason "ai" has no value in the proto Provider enum, which
// lists the providers that have a provider plugin, and providers.ToProto maps
// it to PROVIDER_UNSPECIFIED.
package ai

import (
	"github.com/infracost/go-proto/pkg/tree/ai/code"
)

// AI is the provider node for AI model usage.
type AI struct {
	Code code.Code `tree:"code"`
}
