// Package project holds the canonical definition of an IaC project type: the
// value a parser plugin reports via GetParserConfig's config_file_project_type,
// which the config file records and ProjectInfo.type carries over the wire.
//
// This package is the single source of truth for those values. The config
// package aliases them rather than declaring its own, so config, the CLI and
// the runner all agree on the spelling. It depends only on the standard
// library, so it can sit at the bottom of the dependency graph.
//
// It holds two distinct resolutions, which must not be confused:
//
//   - Resolve answers "which parser plugin parses this project?"
//   - NormalizeForFilter answers "which type family does a governance policy
//     target?"
//
// They differ on types that have their own parser but are not worth surfacing
// as a filter option, such as Cisco Stacks.
package project

import (
	"os"
	"path/filepath"
	"strings"
)

// Type is the project type reported by a parser plugin. It is an open string
// rather than an enum: a plugin may return a custom type via
// config_file_project_type, defaulting to its own name when unset.
type Type string

const (
	// Unknown is a project with no recorded type. Treated as Terraform, which
	// is what the config file assumes for an untyped project.
	Unknown        Type = ""
	Terraform      Type = "terraform"
	Terragrunt     Type = "terragrunt"
	CloudFormation Type = "cloudformation"
	CDKTypeScript  Type = "cdk_typescript"
	CDKJavaScript  Type = "cdk_javascript"
	CDKPython      Type = "cdk_python"
	CiscoStacks    Type = "cisco_stacks"
	Kubernetes     Type = "kubernetes"
	ARM            Type = "arm"
)

// IsCDK reports whether a project is one of the CDK languages.
//
// Matched by prefix rather than against the known CDK constants, so a future
// cdk_<language> is handled without a change here.
func IsCDK(t Type) bool {
	return strings.HasPrefix(string(t), "cdk_")
}

// ResolveUntyped fills in a missing project type by probing its directory.
//
// An Unknown type comes from a config file written before the type was
// recorded, so the directory is probed to tell Terragrunt from Terraform. A
// project that already has a type is returned unchanged.
//
// This deliberately does not fold CDK onto CloudFormation: callers that need
// the project's own type, such as deciding whether to default a CDK project's
// AWS region, must still be able to see it. Compose with ParserFamily, or use
// Resolve, when the parser plugin is what's wanted.
func ResolveUntyped(t Type, dir string) Type {
	if t != Unknown {
		return t
	}

	for _, name := range []string{"terragrunt.hcl", "terragrunt.hcl.json"} {
		if stat, err := os.Stat(filepath.Join(dir, name)); err == nil && !stat.IsDir() {
			return Terragrunt
		}
	}

	return Terraform
}

// ParserFamily returns the type whose parser plugin handles t.
//
// CDK projects are preprocessed into CloudFormation templates and so reuse the
// cloudformation plugin. Everything else passes through unchanged, including
// types like Cisco Stacks that have a parser of their own.
//
// Callers wanting the type for governance filtering want NormalizeForFilter
// instead: this would load the right plugin but filter the wrong family.
func ParserFamily(t Type) Type {
	if IsCDK(t) {
		return CloudFormation
	}

	return t
}

// Resolve returns the project type to use when looking up a parser plugin:
// ResolveUntyped followed by ParserFamily. It is idempotent, so a caller that
// has already resolved a type can safely call it again.
func Resolve(t Type, dir string) Type {
	return ParserFamily(ResolveUntyped(t, dir))
}

// Filterable is the set of types a governance policy can be filtered by, in
// display order. It is deliberately coarser than the full set above:
// NormalizeForFilter folds the niche and derived types onto the family a user
// would recognise, so policies are written against a short, stable list.
var Filterable = []Type{Terraform, CloudFormation, ARM, Kubernetes}

// NormalizeForFilter collapses a project type onto the Filterable set, so a
// policy targeting "cloudformation" also covers the CDK variants that the
// cloudformation plugin parses, and one targeting "terraform" covers the whole
// Terraform family.
//
// The Terraform family is Terragrunt and Cisco Stacks as well as Terraform
// itself: all three are the same HCL tags on the same resources, so a policy
// written for one applies unchanged to the others. Note this means a policy
// cannot target Terragrunt without also covering Terraform.
//
// Anything unrecognised is returned unchanged: a new plugin's type must not
// silently fall into an existing family before it is added above. Callers
// filtering on the result will simply not match it, which is the safe default.
//
// This is not ParserFamily: each of these types has its own parser plugin, and
// folding them there would load the wrong one.
func NormalizeForFilter(t Type) Type {
	switch t {
	case Unknown, Terraform, Terragrunt, CiscoStacks:
		return Terraform
	case CloudFormation, CDKTypeScript, CDKJavaScript, CDKPython:
		return CloudFormation
	case "terraform-plan":
		// Not a declared project type: the terraform-plan plugin sets no
		// config_file_project_type, so the caller falls back to the plugin name.
		return Terraform
	default:
		return t
	}
}
