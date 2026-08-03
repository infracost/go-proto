// Package flag provides bitfield flags for tracking value origin and metadata
// across infrastructure-as-code analysis.
package flag

import (
	"fmt"
	"sort"
	"strings"
)

// Flags is a bitfield for tracking value origin and characteristics.
type Flags uint64

// Flag constants for tracking value sources and characteristics.
const (
	TerraformCode      Flags = 1 << iota // defined in terraform source code
	TerragruntCode                       // defined in terragrunt code
	Synthetic                            // synthesized by us, usually as the value was missing
	Config                               // defined in infracost config file
	EnvVar                               // defined in env var
	TFVars                               // defined in terraform vars e.g. -var or infracost config
	TFVarsFile                           // defined in terraform vars file
	LocalModule                          // originated in a non-root local module
	RemoteModule                         // originated in a remote module
	RegistryModule                       // originated in a registry module
	Sensitive                            // defined as sensitive, e.g. a password
	YorConfig                            // defined in yor config file
	Modified                             // created/modified in a diff
	UsageFile                            // defined in the usage file
	Spacelift                            // defined in spacelift
	TerraformCloud                       // defined in terraform cloud
	CloudFormationYAML                   // defined in cloudformation yaml source code
	CloudFormationJSON                   // defined in cloudformation json source code
	JSONEncoded                          // this variable is in JSON format (e.g. passed by terragrunt)
	TerraformPlanJSON                    // defined in terraform plan JSON output
	Defaulted                            // defaulted based on e.g. cloud api specs, terraform provider specs
	Partial                              // partial/runtime-unknown object: absent keys are unknown state, not errors (e.g. an ARM reference() result)
	// <-- new values here please!

	flagMax // for testing - MUST remain at end of enum!
)

var names = map[Flags]string{
	TerraformCode:      "terraform",
	TerragruntCode:     "terragrunt",
	Config:             "config",
	EnvVar:             "env",
	Synthetic:          "synthetic",
	TFVars:             "tfvars",
	TFVarsFile:         "tfvarsfile",
	LocalModule:        "local_module",
	RemoteModule:       "remote_module",
	RegistryModule:     "registry_module",
	Sensitive:          "sensitive",
	YorConfig:          "yor",
	Modified:           "modified",
	UsageFile:          "usage_file",
	Spacelift:          "spacelift",
	TerraformCloud:     "terraform_cloud",
	CloudFormationYAML: "cloudformation_yaml",
	CloudFormationJSON: "cloudformation_json",
	JSONEncoded:        "json_encoded",
	TerraformPlanJSON:  "terraform_plan_json",
	Defaulted:          "defaulted",
	Partial:            "partial",
}

// DescribeFlags returns a human-readable string describing the set flags.
func DescribeFlags(flags uint64) string {
	var parts []string
	for flag, name := range names {
		if flags&uint64(flag) != 0 {
			parts = append(parts, name)
		}
	}
	sort.Strings(parts)
	return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
}

// IsSynthetic returns true if the Synthetic flag is set.
func (f Flags) IsSynthetic() bool {
	return f&Synthetic != 0
}

// IsSensitive returns true if the Sensitive flag is set.
func (f Flags) IsSensitive() bool {
	return f&Sensitive != 0
}

// IsCode returns true if the TerraformCode flag is set.
func (f Flags) IsCode() bool {
	return f&TerraformCode != 0
}

// IsUnknown returns true if no flags are set.
func (f Flags) IsUnknown() bool {
	return f == 0
}

// IsModule returns true if any module flag (Local, Remote, or Registry) is set.
// Note that this returns false for a root module.
func (f Flags) IsModule() bool {
	return f&(LocalModule|RemoteModule|RegistryModule) != 0
}

// IsLocalModule returns true if the LocalModule flag is set.
func (f Flags) IsLocalModule() bool {
	return f&LocalModule != 0
}

// IsRemoteModule returns true if the RemoteModule flag is set.
func (f Flags) IsRemoteModule() bool {
	return f&RemoteModule != 0
}

// IsRegistryModule returns true if the RegistryModule flag is set.
func (f Flags) IsRegistryModule() bool {
	return f&RegistryModule != 0
}

// IsPurelySynthetic returns true if only the Synthetic flag is set.
func (f Flags) IsPurelySynthetic() bool {
	return f == Synthetic
}

// IsPurelyCode returns true if only the TerraformCode flag is set.
func (f Flags) IsPurelyCode() bool {
	return f == TerraformCode
}

// IsPurelyUnknown returns true if no flags are set (alias for IsUnknown).
func (f Flags) IsPurelyUnknown() bool {
	return f == 0
}

// IsTerragrunt returns true if the TerragruntCode flag is set.
func (f Flags) IsTerragrunt() bool {
	return f&TerragruntCode != 0
}

// IsTerraform returns true if the TerraformCode flag is set.
func (f Flags) IsTerraform() bool {
	return f&TerraformCode != 0
}

// IsCloudFormation returns true if either CloudFormation flag (YAML or JSON) is set.
func (f Flags) IsCloudFormation() bool {
	return f&CloudFormationYAML != 0 || f&CloudFormationJSON != 0
}

// IsCloudFormationYAML returns true if the CloudFormationYAML flag is set.
func (f Flags) IsCloudFormationYAML() bool {
	return f&CloudFormationYAML != 0
}

// IsCloudFormationJSON returns true if the CloudFormationJSON flag is set.
func (f Flags) IsCloudFormationJSON() bool {
	return f&CloudFormationJSON != 0
}

// IsJSONEncoded returns true if the JSONEncoded flag is set.
func (f Flags) IsJSONEncoded() bool {
	return f&JSONEncoded != 0
}

// IsDefault returns true if the Defaulted flag is set
func (f Flags) IsDefault() bool {
	return f&Defaulted != 0
}

// ToProto returns the flags as a uint64 for protocol buffer serialization.
func (f Flags) ToProto() uint64 {
	return uint64(f)
}

// FromProto creates Flags from a protocol buffer uint64 representation.
func FromProto(f uint64) Flags {
	return Flags(f)
}
