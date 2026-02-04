# Flag

Bitfield flags for tracking value origin and metadata.

## Overview

Flags track where values originated (Terraform, Terragrunt, config files, environment variables) and their characteristics (sensitive, modified, module type). Multiple flags can be combined using bitwise OR.

## Usage

```go
import "github.com/infracost/go-proto/pkg/flag"

// Combine flags
f := flag.TerraformCode | flag.LocalModule | flag.Sensitive

// Check flags
if f.IsTerraform() { ... }
if f.IsSensitive() { ... }
if f.IsModule() { ... }

// Describe for debugging
fmt.Println(flag.DescribeFlags(uint64(f))) // [local_module, sensitive, terraform]

// Proto conversion
proto := f.ToProto()
f := flag.FromProto(proto)
```

## Flag Categories

**Source Flags**
- `TerraformCode`, `TerragruntCode` - IaC source code
- `CloudFormationYAML`, `CloudFormationJSON` - CloudFormation templates
- `Config`, `EnvVar`, `TFVars`, `TFVarsFile` - Configuration sources
- `UsageFile`, `Spacelift`, `TerraformCloud` - External sources

**Module Flags**
- `LocalModule`, `RemoteModule`, `RegistryModule`

**Metadata Flags**
- `Synthetic` - Value was synthesized (not from source)
- `Sensitive` - Contains sensitive data
- `Modified` - Created or modified in a diff
- `JSONEncoded` - Value is JSON-encoded
