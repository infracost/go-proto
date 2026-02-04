# Diagnostic

Wrappers and utilities for working with [diagnostic protos](https://github.com/infracost/proto/blob/main/proto/infracost/parser/diagnostic.proto) in Go.

## Overview

Diagnostics represent errors, warnings, and other issues encountered during parsing and analysis. They are categorized by severity (critical, warning, transient) and can include source location information.

## Usage

```go
import "github.com/infracost/go-proto/pkg/diagnostic"

// Create a diagnostic
d := diagnostic.New(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_HCL_PARSE_ERROR, "failed to parse: %s", filename)

// From an error
d := diagnostic.FromError(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_FILESYSTEM_ERROR, err)

// Add metadata
d = d.WithSourceRange(sourceRange).WithLabel("module", "vpc")

// Work with collections
diags := diagnostic.NewDiagnostics(nil)
diags = diags.Add(d1, d2)
diags = diags.Merge(otherDiags)

// Filter diagnostics
critical := diags.Critical()
warnings := diags.NonCritical()
transient := diags.Transient()
byType := diags.OfType(parserpb.DiagnosticType_DIAGNOSTIC_TYPE_HCL_PARSE_ERROR)

// Proto conversion
proto := diags.ToProto()
diags := diagnostic.FromProto(proto)
```

## Severity Categories

- **Critical**: Parse errors, invalid configurations, security issues, cyclic dependencies
- **Warning**: Missing inputs, filesystem errors, unsupported CloudFormation features
- **Transient**: Module fetch errors, source map issues (may resolve on retry)
