// Package ai is the AI provider in the resource tree. It holds where and how
// AI models are used, as opposed to the cloud resources that host them.
//
// The code service holds LLM call sites found in application source code by
// the appcode parser plugin: the SDK clients, the calls made through them, and
// the model IDs those calls use.
//
// AI is a top-level section of the tree alongside aws and azure, filled in by
// a parser plugin like the others. Like them, it has a matching provider
// plugin, which reads only this section and evaluates FinOps policies against
// these resources. Unlike the cloud sections, nothing prices them. Other
// provider plugins are handed the whole tree, but each reads only its own
// section and skips this one.
package ai

import (
	"github.com/infracost/go-proto/pkg/tree/ai/code"
)

// AI is the provider node for AI model usage.
type AI struct {
	Code code.Code `tree:"code"`
}
