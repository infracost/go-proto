package project

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeForFilter(t *testing.T) {
	tests := []struct {
		name     string
		input    Type
		expected Type
	}{
		{"terraform is unchanged", Terraform, Terraform},
		{"terragrunt stays distinct", Terragrunt, Terragrunt},
		{"cloudformation is unchanged", CloudFormation, CloudFormation},
		{"kubernetes is unchanged", Kubernetes, Kubernetes},

		{"untyped project folds to terraform", Unknown, Terraform},
		{"cisco stacks folds to terraform", CiscoStacks, Terraform},
		{"terraform-plan folds to terraform", "terraform-plan", Terraform},

		{"cdk typescript folds to cloudformation", CDKTypeScript, CloudFormation},
		{"cdk javascript folds to cloudformation", CDKJavaScript, CloudFormation},
		{"cdk python folds to cloudformation", CDKPython, CloudFormation},

		{"unrecognised type is left alone", "pulumi", "pulumi"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, NormalizeForFilter(tt.input))
		})
	}
}

// Every type must normalize onto the set users can actually pick from,
// otherwise a policy could never match projects of that type.
func TestNormalizeForFilter_LandsOnFilterableSet(t *testing.T) {
	all := []Type{
		Unknown, Terraform, Terragrunt, CloudFormation,
		CDKTypeScript, CDKJavaScript, CDKPython, CiscoStacks, Kubernetes,
	}

	for _, projectType := range all {
		t.Run(string(projectType), func(t *testing.T) {
			assert.Contains(t, Filterable, NormalizeForFilter(projectType))
		})
	}
}