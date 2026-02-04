# Address

Models Terraform addresses. Wraps the [address proto](https://github.com/infracost/proto/blob/main/proto/infracost/parser/address.proto), adding utility methods.

## Overview

An address represents a path to a resource, variable, or other entity in Terraform configuration. Addresses are composed of segments that can include names, integer indices (e.g., `[0]`), and string indices (e.g., `["key"]`).

Example addresses:
- `module.vpc.aws_subnet.public[0]`
- `local.settings["environment"]`
- `var.instance_count`

This could theoretically be used for other configuration languages as well, but is currently focused on Terraform.

## Usage

```go
import "github.com/infracost/go-proto/pkg/address"

// Create an address
addr := address.New("module", "vpc", "aws_subnet", "public")

// Parse an address string
addr, err := address.Parse("module.vpc.aws_subnet.public[0]")

// Add an index
indexed := addr.CreateIntIndexedChild(0)

// Get the module portion
mod := addr.Module()  // module.vpc

// Get the local (non-module) portion
local := addr.Local() // aws_subnet.public[0]

// Check if one address starts with another
if addr.StartsWith(other) { ... }
```

## Key Types

- `Address` - The main address type wrapping the proto
- `Iterator` - For traversing address segments with type awareness
