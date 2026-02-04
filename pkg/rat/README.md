# Rat

Arbitrary-precision rational numbers. Wraps Go's `math/big.Rat` with convenience methods and proto serialization to `infracost.rat.Rat` from `infracost/proto`.

## Overview

The `Rat` type provides exact rational arithmetic without floating-point precision loss. Useful for financial calculations and other scenarios where precision matters.

## Usage

```go
import "github.com/infracost/go-proto/pkg/rat"

// Create from various numeric types
r := rat.New(42)
r := rat.New(3.14159)
r, err := rat.NewFromString("1/3")

// Arithmetic
sum := a.Add(b)
diff := a.Sub(b)
product := a.Mul(b)
quotient := a.Div(b)

// Comparisons
if a.GreaterThan(b) { ... }
if a.IsZero() { ... }
if a.Equals(b) { ... }

// Rounding
floored := r.Floor()
ceiled := r.Ceil()
rounded := r.Round(2) // 2 decimal places

// Conversions
f := r.Float64()
i := r.Int64()
s := r.String()
s := r.StringFixed(2) // "1.23"

// Proto serialization
proto := r.Proto()
r := rat.FromProto(proto)
```

## Key Types

- `Rat` - Arbitrary-precision rational number
- `Numeric` - Generic constraint for numeric input types
- `Zero` - Pre-defined zero value
