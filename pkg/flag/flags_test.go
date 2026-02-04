package flag

import (
	"testing"
)

func Test_FlagsDescribe(t *testing.T) {
	for i := Flags(1); i != flagMax; i <<= 1 {
		if DescribeFlags(uint64(i)) == "[]" {
			t.Errorf("flag with value %d (0b%b) is not configured for string conversion", i, i)
		}
	}
}

func Test_IsSynthetic(t *testing.T) {
	tests := []struct {
		name     string
		flags    Flags
		expected bool
	}{
		{"synthetic only", Synthetic, true},
		{"synthetic with other", Synthetic | TerraformCode, true},
		{"not synthetic", TerraformCode, false},
		{"zero", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flags.IsSynthetic(); got != tt.expected {
				t.Errorf("IsSynthetic() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_IsSensitive(t *testing.T) {
	tests := []struct {
		name     string
		flags    Flags
		expected bool
	}{
		{"sensitive only", Sensitive, true},
		{"sensitive with other", Sensitive | TerraformCode, true},
		{"not sensitive", TerraformCode, false},
		{"zero", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flags.IsSensitive(); got != tt.expected {
				t.Errorf("IsSensitive() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_IsCode(t *testing.T) {
	tests := []struct {
		name     string
		flags    Flags
		expected bool
	}{
		{"terraform code", TerraformCode, true},
		{"terraform code with other", TerraformCode | Sensitive, true},
		{"not code", Synthetic, false},
		{"zero", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flags.IsCode(); got != tt.expected {
				t.Errorf("IsCode() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_IsUnknown(t *testing.T) {
	tests := []struct {
		name     string
		flags    Flags
		expected bool
	}{
		{"zero", 0, true},
		{"has flag", TerraformCode, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flags.IsUnknown(); got != tt.expected {
				t.Errorf("IsUnknown() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_IsModule(t *testing.T) {
	tests := []struct {
		name     string
		flags    Flags
		expected bool
	}{
		{"local module", LocalModule, true},
		{"remote module", RemoteModule, true},
		{"registry module", RegistryModule, true},
		{"local and remote", LocalModule | RemoteModule, true},
		{"not module", TerraformCode, false},
		{"zero", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flags.IsModule(); got != tt.expected {
				t.Errorf("IsModule() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_IsLocalModule(t *testing.T) {
	tests := []struct {
		name     string
		flags    Flags
		expected bool
	}{
		{"local module", LocalModule, true},
		{"local with other", LocalModule | Sensitive, true},
		{"remote module", RemoteModule, false},
		{"zero", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flags.IsLocalModule(); got != tt.expected {
				t.Errorf("IsLocalModule() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_IsRemoteModule(t *testing.T) {
	tests := []struct {
		name     string
		flags    Flags
		expected bool
	}{
		{"remote module", RemoteModule, true},
		{"remote with other", RemoteModule | Sensitive, true},
		{"local module", LocalModule, false},
		{"zero", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flags.IsRemoteModule(); got != tt.expected {
				t.Errorf("IsRemoteModule() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_IsRegistryModule(t *testing.T) {
	tests := []struct {
		name     string
		flags    Flags
		expected bool
	}{
		{"registry module", RegistryModule, true},
		{"registry with other", RegistryModule | Sensitive, true},
		{"local module", LocalModule, false},
		{"zero", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flags.IsRegistryModule(); got != tt.expected {
				t.Errorf("IsRegistryModule() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_IsPurelySynthetic(t *testing.T) {
	tests := []struct {
		name     string
		flags    Flags
		expected bool
	}{
		{"synthetic only", Synthetic, true},
		{"synthetic with other", Synthetic | TerraformCode, false},
		{"not synthetic", TerraformCode, false},
		{"zero", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flags.IsPurelySynthetic(); got != tt.expected {
				t.Errorf("IsPurelySynthetic() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_IsPurelyCode(t *testing.T) {
	tests := []struct {
		name     string
		flags    Flags
		expected bool
	}{
		{"terraform code only", TerraformCode, true},
		{"terraform code with other", TerraformCode | Sensitive, false},
		{"not code", Synthetic, false},
		{"zero", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flags.IsPurelyCode(); got != tt.expected {
				t.Errorf("IsPurelyCode() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_IsPurelyUnknown(t *testing.T) {
	tests := []struct {
		name     string
		flags    Flags
		expected bool
	}{
		{"zero", 0, true},
		{"has flag", TerraformCode, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flags.IsPurelyUnknown(); got != tt.expected {
				t.Errorf("IsPurelyUnknown() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_IsTerragrunt(t *testing.T) {
	tests := []struct {
		name     string
		flags    Flags
		expected bool
	}{
		{"terragrunt", TerragruntCode, true},
		{"terragrunt with other", TerragruntCode | Sensitive, true},
		{"terraform", TerraformCode, false},
		{"zero", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flags.IsTerragrunt(); got != tt.expected {
				t.Errorf("IsTerragrunt() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_IsTerraform(t *testing.T) {
	tests := []struct {
		name     string
		flags    Flags
		expected bool
	}{
		{"terraform", TerraformCode, true},
		{"terraform with other", TerraformCode | Sensitive, true},
		{"terragrunt", TerragruntCode, false},
		{"zero", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flags.IsTerraform(); got != tt.expected {
				t.Errorf("IsTerraform() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_IsCloudFormation(t *testing.T) {
	tests := []struct {
		name     string
		flags    Flags
		expected bool
	}{
		{"cloudformation yaml", CloudFormationYAML, true},
		{"cloudformation json", CloudFormationJSON, true},
		{"both", CloudFormationYAML | CloudFormationJSON, true},
		{"yaml with other", CloudFormationYAML | Sensitive, true},
		{"terraform", TerraformCode, false},
		{"zero", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flags.IsCloudFormation(); got != tt.expected {
				t.Errorf("IsCloudFormation() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_IsCloudFormationYAML(t *testing.T) {
	tests := []struct {
		name     string
		flags    Flags
		expected bool
	}{
		{"cloudformation yaml", CloudFormationYAML, true},
		{"yaml with other", CloudFormationYAML | Sensitive, true},
		{"cloudformation json", CloudFormationJSON, false},
		{"zero", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flags.IsCloudFormationYAML(); got != tt.expected {
				t.Errorf("IsCloudFormationYAML() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_IsCloudFormationJSON(t *testing.T) {
	tests := []struct {
		name     string
		flags    Flags
		expected bool
	}{
		{"cloudformation json", CloudFormationJSON, true},
		{"json with other", CloudFormationJSON | Sensitive, true},
		{"cloudformation yaml", CloudFormationYAML, false},
		{"zero", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flags.IsCloudFormationJSON(); got != tt.expected {
				t.Errorf("IsCloudFormationJSON() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_ToProto(t *testing.T) {
	tests := []struct {
		name     string
		flags    Flags
		expected uint64
	}{
		{"zero", 0, 0},
		{"single flag", TerraformCode, uint64(TerraformCode)},
		{"multiple flags", TerraformCode | Sensitive, uint64(TerraformCode | Sensitive)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flags.ToProto(); got != tt.expected {
				t.Errorf("ToProto() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_FromProto(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expected Flags
	}{
		{"zero", 0, 0},
		{"single flag", uint64(TerraformCode), TerraformCode},
		{"multiple flags", uint64(TerraformCode | Sensitive), TerraformCode | Sensitive},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromProto(tt.input); got != tt.expected {
				t.Errorf("FromProto() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func Test_DescribeFlags(t *testing.T) {
	tests := []struct {
		name     string
		flags    uint64
		contains []string
	}{
		{"zero", 0, []string{}},
		{"terraform", uint64(TerraformCode), []string{"terraform"}},
		{"multiple", uint64(TerraformCode | Sensitive), []string{"terraform", "sensitive"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DescribeFlags(tt.flags)
			for _, s := range tt.contains {
				if !contains(got, s) {
					t.Errorf("DescribeFlags() = %v, want it to contain %v", got, s)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
