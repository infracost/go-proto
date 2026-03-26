// Package rat provides arbitrary-precision rational numbers wrapping Go's math/big.Rat
// with convenience methods and protocol buffer serialization.
package rat

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/infracost/proto/gen/go/infracost/rational"
)

// Rat represents an arbitrary-precision rational number.
type Rat struct {
	rational *big.Rat
}

// BigRat returns the underlying big.Rat value.
func (r *Rat) BigRat() *big.Rat {
	return r.rational
}

// FromProto creates a Rat from a protocol buffer representation.
func FromProto(r *rational.Rat) *Rat {
	if r == nil {
		return Zero
	}
	if len(r.Denominator) == 0 {
		return Zero // invalid rational, treat as zero
	}
	if len(r.Denominator) == 1 && r.Denominator[0] == 0 {
		return Zero // invalid rational, treat as zero
	}
	num := new(big.Int).SetBytes(r.Numerator)
	if r.Negative {
		num.Neg(num)
	}
	return wrap(new(big.Rat).SetFrac(
		num,
		new(big.Int).SetBytes(r.Denominator),
	))
}

// Proto returns the protocol buffer representation.
func (r *Rat) Proto() *rational.Rat {
	if r == nil || r.rational == nil {
		return nil
	}
	return &rational.Rat{
		Numerator:   r.rational.Num().Bytes(),
		Denominator: r.rational.Denom().Bytes(),
		Negative:    r.rational.Sign() < 0,
	}
}

// MarshalJSON implements json.Marshaler.
func (r *Rat) MarshalJSON() ([]byte, error) {
	if r == nil || r.rational == nil {
		return []byte("null"), nil
	}
	f, _ := r.rational.Float64()
	return json.Marshal(strconv.FormatFloat(f, 'f', -1, 64)) // e.g. "1.23"
}

// UnmarshalJSON implements json.Unmarshaler.
func (r *Rat) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "null" {
		r.rational = nil
		return nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	r.rational = new(big.Rat).SetFloat64(f)
	return nil
}

// Numeric is a type constraint for numeric types that can be converted to Rat.
type Numeric interface {
	~int | ~int16 | ~int32 | ~int64 |
		~uint | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 | *big.Int |
		*int | *int16 | *int32 | *int64 |
		*uint | *uint16 | *uint32 | *uint64 |
		*float32 | *float64
}

// NewFromString parses a string as a rational number (e.g., "1/3" or "3.14").
func NewFromString(s string) (*Rat, error) {
	b, ok := new(big.Rat).SetString(s)
	if !ok {
		return nil, fmt.Errorf("cannot convert %q to rational number", s)
	}
	return wrap(b), nil
}

// New creates a Rat from any numeric type.
func New[T Numeric](num T) *Rat {
	b := new(big.Rat)
	switch v := any(num).(type) {
	case *int:
		if v == nil {
			return nil
		}
		b.SetInt64(int64(*v))
	case *int16:
		if v == nil {
			return nil
		}
		b.SetInt64(int64(*v))
	case *int32:
		if v == nil {
			return nil
		}
		b.SetInt64(int64(*v))
	case *int64:
		if v == nil {
			return nil
		}
		b.SetInt64(*v)
	case *uint:
		if v == nil {
			return nil
		}
		b.SetUint64(uint64(*v))
	case *uint16:
		if v == nil {
			return nil
		}
		b.SetUint64(uint64(*v))
	case *uint32:
		if v == nil {
			return nil
		}
		b.SetUint64(uint64(*v))
	case *uint64:
		if v == nil {
			return nil
		}
		b.SetUint64(*v)
	case *float32:
		if v == nil {
			return nil
		}
		b.SetFloat64(float64(*v))
	case *float64:
		if v == nil {
			return nil
		}
		b.SetFloat64(*v)
	case int:
		b.SetInt64(int64(v))
	case int16:
		b.SetInt64(int64(v))
	case int32:
		b.SetInt64(int64(v))
	case int64:
		b.SetInt64(v)
	case uint:
		b.SetUint64(uint64(v))
	case uint16:
		b.SetUint64(uint64(v))
	case uint32:
		b.SetUint64(uint64(v))
	case uint64:
		b.SetUint64(v)
	case float32:
		b.SetFloat64(float64(v))
	case float64:
		b.SetFloat64(v)
	case *big.Int:
		b.SetInt(v)
	}
	return &Rat{
		rational: b,
	}
}

func wrap(b *big.Rat) *Rat {
	if b == nil {
		return nil
	}
	return &Rat{
		rational: b,
	}
}

// Int64 returns the integer part as int64 (truncates toward zero).
func (r *Rat) Int64() int64 {
	if r == nil || r.rational == nil {
		return 0
	}
	// Use big.Int.Div to avoid int64 overflow of the denominator causing divide by zero.
	result := new(big.Int).Div(r.rational.Num(), r.rational.Denom())
	return result.Int64()
}

// Int returns the integer part as int (truncates toward zero).
func (r *Rat) Int() int {
	if r == nil || r.rational == nil {
		return 0
	}
	return int(r.Int64())
}

// Float64 returns the value as float64 (may lose precision).
func (r *Rat) Float64() float64 {
	if r == nil || r.rational == nil {
		return 0
	}
	f, _ := r.rational.Float64()
	return f
}

// Max returns the greater of r and b.
func (r *Rat) Max(b *Rat) *Rat {
	if r.GreaterThanOrEqual(b) {
		return r
	}
	return b
}

// Min returns the lesser of r and b.
func (r *Rat) Min(b *Rat) *Rat {
	if r.LessThanOrEqual(b) {
		return r
	}
	return b
}

// Mul returns r * b.
func (r *Rat) Mul(b *Rat) *Rat {
	if r == nil || b == nil || r.rational == nil || b.rational == nil {
		return nil
	}
	return wrap(new(big.Rat).Mul(r.rational, b.rational))
}

// Div returns r / b. Returns nil if b is zero or nil.
func (r *Rat) Div(b *Rat) *Rat {
	if r == nil || b == nil || r.rational == nil || b.rational == nil {
		return nil
	}
	if b.rational.Sign() == 0 {
		return nil
	}
	return wrap(new(big.Rat).Quo(r.rational, b.rational))
}

// Add returns r + b.
func (r *Rat) Add(b *Rat) *Rat {
	if r == nil || r.rational == nil {
		return b
	}
	if b == nil || b.rational == nil {
		return r
	}
	return wrap(new(big.Rat).Add(r.rational, b.rational))
}

// Sub returns r - b.
func (r *Rat) Sub(b *Rat) *Rat {
	if r == nil || r.rational == nil {
		return b.Neg()
	}
	if b == nil || b.rational == nil {
		return r
	}
	return wrap(new(big.Rat).Sub(r.rational, b.rational))
}

// String returns the decimal string representation.
func (r *Rat) String() string {
	if r == nil || r.rational == nil {
		return ""
	}
	f, _ := r.rational.Float64()
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// StringFixed returns the decimal string with exactly the given number of decimal places.
func (r *Rat) StringFixed(precision int) string {
	if r == nil || r.rational == nil {
		return ""
	}
	return r.rational.FloatString(precision)
}

// GreaterThan returns true if r > b.
func (r *Rat) GreaterThan(b *Rat) bool {
	if r == nil || r.rational == nil {
		return Zero.GreaterThan(b)
	}
	if b == nil || b.rational == nil {
		return r.GreaterThan(Zero)
	}
	return r.rational.Cmp(b.rational) > 0
}

// Zero is a pre-allocated zero value.
var Zero = New(0)

// GreaterThanZero returns true if r > 0.
func (r *Rat) GreaterThanZero() bool {
	if r == nil || r.rational == nil {
		return false
	}
	return r.rational.Cmp(Zero.rational) > 0
}

// LessThan returns true if r < b.
func (r *Rat) LessThan(b *Rat) bool {
	if r == nil || r.rational == nil {
		return Zero.LessThan(b)
	}
	if b == nil || b.rational == nil {
		return r.LessThan(Zero)
	}
	return r.rational.Cmp(b.rational) < 0
}

// LessThanOrEqual returns true if r <= b.
func (r *Rat) LessThanOrEqual(b *Rat) bool {
	if r == nil || r.rational == nil {
		return Zero.LessThanOrEqual(b)
	}
	if b == nil || b.rational == nil {
		return r.LessThanOrEqual(Zero)
	}
	return r.rational.Cmp(b.rational) <= 0
}

// GreaterThanOrEqual returns true if r >= b.
func (r *Rat) GreaterThanOrEqual(b *Rat) bool {
	if r == nil || r.rational == nil {
		return Zero.GreaterThanOrEqual(b)
	}
	if b == nil || b.rational == nil {
		return r.GreaterThanOrEqual(Zero)
	}
	return r.rational.Cmp(b.rational) >= 0
}

// IsZero returns true if r equals zero or is nil.
func (r *Rat) IsZero() bool {
	return r == nil || r.rational == nil || r.Equals(Zero)
}

// Equals returns true if r == b.
func (r *Rat) Equals(b *Rat) bool {
	if r == nil || r.rational == nil {
		return b == nil || b.IsZero()
	}
	if b == nil || b.rational == nil {
		return r.IsZero()
	}
	return r.rational.Cmp(b.rational) == 0
}

// IntPart returns the integer part of r as a Rat.
func (r *Rat) IntPart() *Rat {
	if r == nil || r.rational == nil {
		return nil
	}
	return wrap(new(big.Rat).SetInt(new(big.Int).Div(r.rational.Num(), r.rational.Denom())))
}

// Ceil returns the smallest integer >= r.
func (r *Rat) Ceil() *Rat {
	if r == nil || r.rational == nil {
		return nil
	}
	if r.rational.IsInt() {
		return wrap(new(big.Rat).Set(r.rational))
	}
	// Div truncates toward -infinity, add 1 to get ceiling
	i := new(big.Int).Div(r.rational.Num(), r.rational.Denom())
	i.Add(i, big.NewInt(1))
	return wrap(new(big.Rat).SetInt(i))
}

// Floor returns the largest integer <= r.
func (r *Rat) Floor() *Rat {
	if r == nil || r.rational == nil {
		return nil
	}
	if r.rational.IsInt() {
		return wrap(new(big.Rat).Set(r.rational))
	}
	// Div truncates toward -infinity, which is floor behavior
	return wrap(new(big.Rat).SetInt(new(big.Int).Div(r.rational.Num(), r.rational.Denom())))
}

// Round returns r rounded to the given number of decimal places.
func (r *Rat) Round(places int) *Rat {
	if r == nil || r.rational == nil {
		return nil
	}
	if places < 0 {
		places = 0
	}

	// multiplier = 10^places
	mul := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(places)), nil)
	mulRat := new(big.Rat).SetInt(mul)

	// multiply, round to nearest integer, divide
	scaled := new(big.Rat).Mul(r.rational, mulRat)

	// Round to nearest integer
	num := scaled.Num()
	denom := scaled.Denom()

	// integer division with rounding
	quo := new(big.Int).Div(num, denom)
	rem := new(big.Int).Mod(num, denom)

	// if remainder >= denom/2, round up (away from zero)
	rem.Mul(rem, big.NewInt(2))
	if rem.CmpAbs(denom) >= 0 {
		if scaled.Sign() > 0 {
			quo.Add(quo, big.NewInt(1))
		} else {
			quo.Sub(quo, big.NewInt(1))
		}
	}

	// divide back
	result := new(big.Rat).SetInt(quo)
	result.Quo(result, mulRat)

	return wrap(result)
}

// Abs returns the absolute value of r.
func (r *Rat) Abs() *Rat {
	if r == nil || r.rational == nil {
		return nil
	}
	return wrap(new(big.Rat).Abs(r.rational))
}

// ScaledInt converts a Rat to a scaled integer with the given decimal places.
// For example, with decimalPlaces=2: 12.34 becomes 1234
func (r *Rat) ScaledInt(decimalPlaces int) int64 {
	if r == nil || r.rational == nil {
		return 0
	}
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimalPlaces)), nil)
	scaled := new(big.Rat).Mul(r.rational, new(big.Rat).SetInt(scale))

	// Truncate to integer (or use rounding logic from Round() if you prefer)
	return new(big.Int).Div(scaled.Num(), scaled.Denom()).Int64()
}

// Neg returns -r.
func (r *Rat) Neg() *Rat {
	if r == nil || r.rational == nil {
		return nil
	}
	return wrap(new(big.Rat).Neg(r.rational))
}
