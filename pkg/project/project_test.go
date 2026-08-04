package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolve(t *testing.T) {
	tests := []struct {
		name     string
		input    Type
		files    []string
		expected Type
	}{
		{"terraform is unchanged", Terraform, nil, Terraform},
		{"terragrunt is unchanged", Terragrunt, nil, Terragrunt},
		{"cloudformation is unchanged", CloudFormation, nil, CloudFormation},
		{"kubernetes is unchanged", Kubernetes, nil, Kubernetes},

		// Cisco Stacks has its own parser plugin, so unlike filtering this must
		// not fold it onto terraform.
		{"cisco stacks keeps its own parser", CiscoStacks, nil, CiscoStacks},

		{"cdk typescript uses the cloudformation parser", CDKTypeScript, nil, CloudFormation},
		{"cdk javascript uses the cloudformation parser", CDKJavaScript, nil, CloudFormation},
		{"cdk python uses the cloudformation parser", CDKPython, nil, CloudFormation},
		{"an unknown cdk language still uses cloudformation", "cdk_go", nil, CloudFormation},

		{"untyped project defaults to terraform", Unknown, nil, Terraform},
		{"untyped project with terragrunt.hcl is terragrunt", Unknown, []string{"terragrunt.hcl"}, Terragrunt},
		{"untyped project with terragrunt.hcl.json is terragrunt", Unknown, []string{"terragrunt.hcl.json"}, Terragrunt},
		{"untyped project with only main.tf is terraform", Unknown, []string{"main.tf"}, Terraform},

		{"unrecognised type is left alone", "pulumi", nil, "pulumi"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			for _, name := range tt.files {
				require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte("{}"), 0o600))
			}

			assert.Equal(t, tt.expected, Resolve(tt.input, dir))
		})
	}
}

// A directory named terragrunt.hcl must not be mistaken for the file.
func TestResolve_IgnoresTerragruntDirectory(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(dir, "terragrunt.hcl"), 0o750))

	assert.Equal(t, Terraform, Resolve(Unknown, dir))
}

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
