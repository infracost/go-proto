// Package ai is the AI provider in the resource tree. It holds where and how
// AI models are used, as opposed to the cloud resources that host them.
//
// The code service holds LLM call sites found in application source code by
// the appcode parser plugin: the SDK clients, the calls made through them, and
// the model IDs those calls use.
package ai

import (
	"github.com/infracost/go-proto/pkg/tree/ai/code"
)

// AI is the provider node for AI model usage.
type AI struct {
	Code code.Code `tree:"code"`
}
