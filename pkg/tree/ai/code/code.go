// Package code holds LLM usage found in application source code.
//
// A ModelCall is one call site. It references the Client it is made through,
// which says where the request goes (direct to the vendor, through a cloud
// provider, or through a proxy), and the ModelConfig its model ID was read
// from when that ID was set away from the call.
//
// Settings read from source code have three states, not two. A call either
// sets a value, does not set it, or sets it in a way the parser cannot read,
// e.g. a Bedrock invoke_model call whose settings sit inside a serialized
// body. Absence is not representable on its own, because an unset value
// serializes as a zero, so the states are carried as follows:
//
//   - set        the value, with a source range pointing at where it was set
//   - not set    the zero value, with no flags
//   - unknown    the zero value, with flag.Synthetic
//
// The same convention is used for core.LimitRangeItem in the Kubernetes
// provider, where flag.Synthetic likewise separates unset from zero.
// Consumers must check IsSynthetic before treating a zero as "not set": an
// unknown max_tokens is not the same as an unbounded one.
package code

type Code struct {
	Clients      []Client      `tree:"clients"`
	ModelCalls   []ModelCall   `tree:"model_calls"`
	ModelConfigs []ModelConfig `tree:"model_configs"`
}

func (c *Code) PostProcess() {
	// Reset relationships this method writes to so PostProcess is idempotent.
	for i := range c.ModelCalls {
		c.ModelCalls[i].Relationships = ModelCallRelationships{}
	}

	for i, call := range c.ModelCalls {
		if !call.ClientID.IsEmpty() {
			for j := range c.Clients {
				if call.ClientID.Equal(c.Clients[j].ID) {
					c.ModelCalls[i].Relationships.Client = &c.Clients[j]
					break
				}
			}
		}
		if !call.ModelConfigID.IsEmpty() {
			for j := range c.ModelConfigs {
				if call.ModelConfigID.Equal(c.ModelConfigs[j].ID) {
					c.ModelCalls[i].Relationships.ModelConfig = &c.ModelConfigs[j]
					break
				}
			}
		}
	}
}
