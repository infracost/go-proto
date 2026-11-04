package code

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// ModelConfig is a model ID set away from the call that uses it: a constant,
// a variable, an env var default, a key in a config file or a deployment.
// Changing the model for every call that reads it means changing it here, not
// at the calls.
type ModelConfig struct {
	resource.Resource `tree:"-"`

	// Kind is where the model ID is set: "const", "variable", "env",
	// "config" or "deployment" (an Azure OpenAI deployment named in the URL
	// an endpoint env var holds).
	Kind value.String `tree:"kind"`

	// Key is the name it is set under: the constant's or variable's name, the
	// env var's name, or the config key.
	Key value.String `tree:"key"`

	// Value is the model ID, with a source range pointing at where it is set.
	Value value.String `tree:"value"`
}
