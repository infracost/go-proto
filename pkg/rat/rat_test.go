package rat_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/infracost/go-proto/pkg/rat"
)

func TestNewRat(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected float64
	}{
		{"int", 42, 42},
		{"int64", int64(100), 100},
		{"float64", 3.14, 3.14},
		{"negative", -5, -5},
		{"zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r *rat.Rat
			switch v := tt.input.(type) {
			case int:
				r = rat.New(v)
			case int64:
				r = rat.New(v)
			case float64:
				r = rat.New(v)
			}
			if r.Float64() != tt.expected {
				t.Errorf("NewRat(%v) = %v, want %v", tt.input, r.Float64(), tt.expected)
			}
		})
	}
}

func TestRatProtoPersistsNegative(t *testing.T) {
	original := rat.New(-42)
	proto := original.Proto()
	restored := rat.FromProto(proto)

	if !original.Equals(restored) {
		t.Errorf("Expected restored Rat to equal original. got %v, want %v", restored.String(), original.String())
	}
}

func TestNewRatNilPointers(t *testing.T) {
	var intPtr *int
	r := rat.New(intPtr)
	if r != nil {
		t.Error("NewRat(nil *int) should return nil")
	}

	var floatPtr *float64
	r = rat.New(floatPtr)
	if r != nil {
		t.Error("NewRat(nil *float64) should return nil")
	}
}

func TestNewRatFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		wantErr  bool
	}{
		{"123/10", 12.3, false},
		{"1.23", 1.23, false},
		{"5", 5, false},
		{"-7/2", -3.5, false},
		{"invalid", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			r, err := rat.NewFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRatFromString(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err == nil && r.Float64() != tt.expected {
				t.Errorf("NewRatFromString(%q) = %v, want %v", tt.input, r.Float64(), tt.expected)
			}
		})
	}
}

func TestArithmetic(t *testing.T) {
	a := rat.New(10)
	b := rat.New(3)

	if got := a.Add(b).Int(); got != 13 {
		t.Errorf("10 + 3 = %d, want 13", got)
	}

	if got := a.Sub(b).Int(); got != 7 {
		t.Errorf("10 - 3 = %d, want 7", got)
	}

	if got := a.Mul(b).Int(); got != 30 {
		t.Errorf("10 * 3 = %d, want 30", got)
	}

	if got := a.Div(b).Float64(); got != 10.0/3.0 {
		t.Errorf("10 / 3 = %v, want %v", got, 10.0/3.0)
	}
}

func TestComparisons(t *testing.T) {
	a := rat.New(5)
	b := rat.New(10)
	c := rat.New(5)

	if !b.GreaterThan(a) {
		t.Error("10 should be greater than 5")
	}

	if !a.LessThan(b) {
		t.Error("5 should be less than 10")
	}

	if !a.Equals(c) {
		t.Error("5 should equal 5")
	}

	if !a.LessThanOrEqual(c) {
		t.Error("5 should be <= 5")
	}

	if !a.GreaterThanOrEqual(c) {
		t.Error("5 should be >= 5")
	}

	if !b.GreaterThanZero() {
		t.Error("10 should be greater than zero")
	}

	if rat.New(-1).GreaterThanZero() {
		t.Error("-1 should not be greater than zero")
	}
}

func TestMinMax(t *testing.T) {
	a := rat.New(5)
	b := rat.New(10)

	if a.Max(b) != b {
		t.Error("Max(5, 10) should be 10")
	}

	if a.Min(b) != a {
		t.Error("Min(5, 10) should be 5")
	}
}

func TestIsZero(t *testing.T) {
	if !rat.New(0).IsZero() {
		t.Error("0 should be zero")
	}

	if rat.New(1).IsZero() {
		t.Error("1 should not be zero")
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
	}

	for _, tt := range tests {
		r := rat.New(tt.input)
		if got := r.Abs().Int(); got != tt.expected {
			t.Errorf("Abs(%d) = %d, want %d", tt.input, got, tt.expected)
		}
	}
}

func TestCeil(t *testing.T) {
	tests := []struct {
		num, denom string
		expected   int64
	}{
		{"5", "1", 5},
		{"5", "2", 3},   // 2.5 -> 3
		{"7", "3", 3},   // 2.33 -> 3
		{"-5", "2", -2}, // -2.5 -> -2 (toward +infinity)
		{"-7", "3", -2}, // -2.33 -> -2 (toward +infinity)
		{"4", "1", 4},
	}

	for _, tt := range tests {
		r, _ := rat.NewFromString(tt.num + "/" + tt.denom)
		if got := r.Ceil().Int64(); got != tt.expected {
			t.Errorf("Ceil(%s/%s) = %d, want %d", tt.num, tt.denom, got, tt.expected)
		}
	}
}

func TestFloor(t *testing.T) {
	tests := []struct {
		num, denom string
		expected   int64
	}{
		{"5", "1", 5},
		{"5", "2", 2},   // 2.5 -> 2
		{"7", "3", 2},   // 2.33 -> 2
		{"-5", "2", -3}, // -2.5 -> -3 (toward -infinity)
		{"-7", "3", -3}, // -2.33 -> -3 (toward -infinity)
	}

	for _, tt := range tests {
		r, _ := rat.NewFromString(tt.num + "/" + tt.denom)
		if got := r.Floor().Int64(); got != tt.expected {
			t.Errorf("Floor(%s/%s) = %d, want %d", tt.num, tt.denom, got, tt.expected)
		}
	}
}

func TestRound(t *testing.T) {
	tests := []struct {
		input    float64
		places   int
		expected string
	}{
		{1.234, 2, "1.23"},
		{1.235, 2, "1.24"},
		{1.999, 2, "2"},
		{-1.235, 2, "-1.24"},
		{1.5, 0, "2"},
	}

	for _, tt := range tests {
		r := rat.New(tt.input)
		if got := r.Round(tt.places).String(); got != tt.expected {
			t.Errorf("Round(%v, %d) = %s, want %s", tt.input, tt.places, got, tt.expected)
		}
	}
}

func TestIntPart(t *testing.T) {
	r, _ := rat.NewFromString("7/3") // 2.33...
	if got := r.IntPart().Int(); got != 2 {
		t.Errorf("IntPart(7/3) = %d, want 2", got)
	}
}

func TestScaledInt(t *testing.T) {
	tests := []struct {
		input    string
		places   int
		expected int64
	}{
		{"1234/100", 2, 1234}, // 12.34
		{"3/2", 1, 15},        // 1.5
		{"100/1", 2, 10000},   // 100.0
		{"1/20", 2, 5},        // 0.05
	}

	for _, tt := range tests {
		r, _ := rat.NewFromString(tt.input)
		if got := r.ScaledInt(tt.places); got != tt.expected {
			t.Errorf("ScaledInt(%s, %d) = %d, want %d", tt.input, tt.places, got, tt.expected)
		}
	}
}

func TestJSONMarshal(t *testing.T) {
	r := rat.New(1.5)
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	if string(b) != `"1.5"` {
		t.Errorf("Marshal = %s, want %s", string(b), `"1.5"`)
	}

	// Test nil
	var nilRat *rat.Rat
	b, err = json.Marshal(nilRat)
	if err != nil {
		t.Fatalf("Marshal nil error: %v", err)
	}
	if string(b) != "null" {
		t.Errorf("Marshal nil = %s, want null", string(b))
	}
}

func TestJSONUnmarshal(t *testing.T) {
	var r rat.Rat
	err := json.Unmarshal([]byte(`"12.3"`), &r)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if r.Float64() != 12.3 {
		t.Errorf("Unmarshal = %v, want 12.3", r.Float64())
	}

	// Test invalid
	err = json.Unmarshal([]byte(`"invalid"`), &r)
	if err == nil {
		t.Error("Unmarshal invalid should error")
	}
}

func TestBigRat(t *testing.T) {
	r := rat.New(42)
	br := r.BigRat()
	if br.Cmp(big.NewRat(42, 1)) != 0 {
		t.Error("BigRat() should return underlying *big.Rat")
	}
}

func TestString(t *testing.T) {
	r, _ := rat.NewFromString("1/3")
	s := r.String()
	if s != "0.3333333333333333" {
		t.Errorf("String() = %s, want 0.3333333333333333", s)
	}
}

func TestProtoRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"positive integer", "42/1"},
		{"positive fraction", "22/7"},
		{"one", "1/1"},
		{"zero", "0/1"},
		{"large numerator", "123456789/1"},
		{"large denominator", "1/123456789"},
		{"large both", "987654321/123456789"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original, err := rat.NewFromString(tt.input)
			if err != nil {
				t.Fatalf("NewFromString(%q) error: %v", tt.input, err)
			}

			proto := original.Proto()
			if proto == nil {
				t.Fatal("Proto() returned nil")
			}

			restored := rat.FromProto(proto)
			if restored == nil {
				t.Fatal("FromProto() returned nil")
			}

			if !original.Equals(restored) {
				t.Errorf("round trip failed: original=%s, restored=%s", original.String(), restored.String())
			}
		})
	}
}

func TestFromProtoNil(t *testing.T) {
	r := rat.FromProto(nil)
	if !r.IsZero() {
		t.Error("FromProto(nil) should return Zero")
	}
}

func TestProtoNil(t *testing.T) {
	var r *rat.Rat
	if r.Proto() != nil {
		t.Error("nil.Proto() should return nil")
	}
}

func TestNilRatOperations(t *testing.T) {
	var r *rat.Rat

	// Test nil receiver returns safe defaults
	if r.Int64() != 0 {
		t.Error("nil.Int64() should return 0")
	}
	if r.Int() != 0 {
		t.Error("nil.Int() should return 0")
	}
	if r.Float64() != 0 {
		t.Error("nil.Float64() should return 0")
	}
	if r.String() != "" {
		t.Error("nil.String() should return empty string")
	}
	if r.StringFixed(2) != "" {
		t.Error("nil.StringFixed() should return empty string")
	}
	if !r.IsZero() {
		t.Error("nil.IsZero() should return true")
	}
	if r.IntPart() != nil {
		t.Error("nil.IntPart() should return nil")
	}
	if r.Ceil() != nil {
		t.Error("nil.Ceil() should return nil")
	}
	if r.Floor() != nil {
		t.Error("nil.Floor() should return nil")
	}
	if r.Round(2) != nil {
		t.Error("nil.Round() should return nil")
	}
	if r.Abs() != nil {
		t.Error("nil.Abs() should return nil")
	}
	if r.Neg() != nil {
		t.Error("nil.Neg() should return nil")
	}
	if r.ScaledInt(2) != 0 {
		t.Error("nil.ScaledInt() should return 0")
	}
	if r.GreaterThanZero() {
		t.Error("nil.GreaterThanZero() should return false")
	}
}

func TestNilArithmetic(t *testing.T) {
	var r *rat.Rat
	b := rat.New(5)

	// Mul with nil
	if r.Mul(b) != nil {
		t.Error("nil.Mul(b) should return nil")
	}
	if b.Mul(r) != nil {
		t.Error("b.Mul(nil) should return nil")
	}

	// Div with nil
	if r.Div(b) != nil {
		t.Error("nil.Div(b) should return nil")
	}
	if b.Div(r) != nil {
		t.Error("b.Div(nil) should return nil")
	}

	// Add with nil
	result := r.Add(b)
	if !result.Equals(b) {
		t.Error("nil.Add(b) should return b")
	}
	result = b.Add(r)
	if !result.Equals(b) {
		t.Error("b.Add(nil) should return b")
	}

	// Sub with nil
	result = r.Sub(b)
	if !result.Equals(b.Neg()) {
		t.Error("nil.Sub(b) should return -b")
	}
	result = b.Sub(r)
	if !result.Equals(b) {
		t.Error("b.Sub(nil) should return b")
	}
}

func TestNilComparisons(t *testing.T) {
	var r *rat.Rat
	b := rat.New(5)
	neg := rat.New(-5)

	// GreaterThan
	if r.GreaterThan(neg) != rat.Zero.GreaterThan(neg) {
		t.Error("nil.GreaterThan should behave like Zero.GreaterThan")
	}
	if b.GreaterThan(r) != b.GreaterThan(rat.Zero) {
		t.Error("b.GreaterThan(nil) should behave like b.GreaterThan(Zero)")
	}

	// LessThan
	if r.LessThan(b) != rat.Zero.LessThan(b) {
		t.Error("nil.LessThan should behave like Zero.LessThan")
	}
	if neg.LessThan(r) != neg.LessThan(rat.Zero) {
		t.Error("neg.LessThan(nil) should behave like neg.LessThan(Zero)")
	}

	// LessThanOrEqual
	if r.LessThanOrEqual(b) != rat.Zero.LessThanOrEqual(b) {
		t.Error("nil.LessThanOrEqual should behave like Zero.LessThanOrEqual")
	}
	if b.LessThanOrEqual(r) != b.LessThanOrEqual(rat.Zero) {
		t.Error("b.LessThanOrEqual(nil) should behave like b.LessThanOrEqual(Zero)")
	}

	// GreaterThanOrEqual
	if r.GreaterThanOrEqual(neg) != rat.Zero.GreaterThanOrEqual(neg) {
		t.Error("nil.GreaterThanOrEqual should behave like Zero.GreaterThanOrEqual")
	}
	if b.GreaterThanOrEqual(r) != b.GreaterThanOrEqual(rat.Zero) {
		t.Error("b.GreaterThanOrEqual(nil) should behave like b.GreaterThanOrEqual(Zero)")
	}
}

func TestNilEquals(t *testing.T) {
	var r *rat.Rat
	var s *rat.Rat

	if !r.Equals(s) {
		t.Error("nil.Equals(nil) should return true")
	}
	if !r.Equals(rat.Zero) {
		t.Error("nil.Equals(Zero) should return true")
	}
	if !rat.Zero.Equals(r) {
		t.Error("Zero.Equals(nil) should return true")
	}
	if r.Equals(rat.New(1)) {
		t.Error("nil.Equals(1) should return false")
	}
}

func TestNeg(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{5, -5},
		{-5, 5},
		{0, 0},
	}

	for _, tt := range tests {
		r := rat.New(tt.input)
		if got := r.Neg().Int(); got != tt.expected {
			t.Errorf("Neg(%d) = %d, want %d", tt.input, got, tt.expected)
		}
	}
}

func TestNewWithPointers(t *testing.T) {
	// int pointers
	intVal := 42
	r := rat.New(&intVal)
	if r.Int() != 42 {
		t.Errorf("New(&int) = %d, want 42", r.Int())
	}

	int16Val := int16(16)
	r = rat.New(&int16Val)
	if r.Int() != 16 {
		t.Errorf("New(&int16) = %d, want 16", r.Int())
	}

	int32Val := int32(32)
	r = rat.New(&int32Val)
	if r.Int() != 32 {
		t.Errorf("New(&int32) = %d, want 32", r.Int())
	}

	int64Val := int64(64)
	r = rat.New(&int64Val)
	if r.Int() != 64 {
		t.Errorf("New(&int64) = %d, want 64", r.Int())
	}

	// uint pointers
	uintVal := uint(100)
	r = rat.New(&uintVal)
	if r.Int() != 100 {
		t.Errorf("New(&uint) = %d, want 100", r.Int())
	}

	uint16Val := uint16(116)
	r = rat.New(&uint16Val)
	if r.Int() != 116 {
		t.Errorf("New(&uint16) = %d, want 116", r.Int())
	}

	uint32Val := uint32(132)
	r = rat.New(&uint32Val)
	if r.Int() != 132 {
		t.Errorf("New(&uint32) = %d, want 132", r.Int())
	}

	uint64Val := uint64(164)
	r = rat.New(&uint64Val)
	if r.Int() != 164 {
		t.Errorf("New(&uint64) = %d, want 164", r.Int())
	}

	// float pointers
	float32Val := float32(3.14)
	r = rat.New(&float32Val)
	if r.Float64() < 3.13 || r.Float64() > 3.15 {
		t.Errorf("New(&float32) = %f, want ~3.14", r.Float64())
	}

	float64Val := 2.718
	r = rat.New(&float64Val)
	if r.Float64() != 2.718 {
		t.Errorf("New(&float64) = %f, want 2.718", r.Float64())
	}
}

func TestNewWithNilPointers(t *testing.T) {
	var intPtr *int
	if rat.New(intPtr) != nil {
		t.Error("New(nil *int) should return nil")
	}

	var int16Ptr *int16
	if rat.New(int16Ptr) != nil {
		t.Error("New(nil *int16) should return nil")
	}

	var int32Ptr *int32
	if rat.New(int32Ptr) != nil {
		t.Error("New(nil *int32) should return nil")
	}

	var int64Ptr *int64
	if rat.New(int64Ptr) != nil {
		t.Error("New(nil *int64) should return nil")
	}

	var uintPtr *uint
	if rat.New(uintPtr) != nil {
		t.Error("New(nil *uint) should return nil")
	}

	var uint16Ptr *uint16
	if rat.New(uint16Ptr) != nil {
		t.Error("New(nil *uint16) should return nil")
	}

	var uint32Ptr *uint32
	if rat.New(uint32Ptr) != nil {
		t.Error("New(nil *uint32) should return nil")
	}

	var uint64Ptr *uint64
	if rat.New(uint64Ptr) != nil {
		t.Error("New(nil *uint64) should return nil")
	}

	var float32Ptr *float32
	if rat.New(float32Ptr) != nil {
		t.Error("New(nil *float32) should return nil")
	}

	var float64Ptr *float64
	if rat.New(float64Ptr) != nil {
		t.Error("New(nil *float64) should return nil")
	}
}

func TestNewWithDirectTypes(t *testing.T) {
	// Direct int types
	if r := rat.New(int16(16)); r.Int() != 16 {
		t.Errorf("New(int16) = %d, want 16", r.Int())
	}
	if r := rat.New(int32(32)); r.Int() != 32 {
		t.Errorf("New(int32) = %d, want 32", r.Int())
	}
	if r := rat.New(uint(100)); r.Int() != 100 {
		t.Errorf("New(uint) = %d, want 100", r.Int())
	}
	if r := rat.New(uint16(116)); r.Int() != 116 {
		t.Errorf("New(uint16) = %d, want 116", r.Int())
	}
	if r := rat.New(uint32(132)); r.Int() != 132 {
		t.Errorf("New(uint32) = %d, want 132", r.Int())
	}
	if r := rat.New(uint64(164)); r.Int() != 164 {
		t.Errorf("New(uint64) = %d, want 164", r.Int())
	}
	if r := rat.New(float32(3.0)); r.Float64() != 3.0 {
		t.Errorf("New(float32) = %f, want 3.0", r.Float64())
	}
}

func TestNewWithBigInt(t *testing.T) {
	bigInt := big.NewInt(999)
	r := rat.New(bigInt)
	if r.Int() != 999 {
		t.Errorf("New(*big.Int) = %d, want 999", r.Int())
	}
}

func TestRoundNegativePlaces(t *testing.T) {
	r := rat.New(1.234)
	// Negative places should be treated as 0
	result := r.Round(-5)
	if result.Int() != 1 {
		t.Errorf("Round with negative places should round to integer, got %s", result.String())
	}
}

func TestStringFixed(t *testing.T) {
	tests := []struct {
		input     float64
		precision int
		expected  string
	}{
		{1.234, 2, "1.23"},
		{1.235, 2, "1.24"},
		{100.0, 2, "100.00"},
		{0.0, 2, "0.00"},
	}

	for _, tt := range tests {
		r := rat.New(tt.input)
		if got := r.StringFixed(tt.precision); got != tt.expected {
			t.Errorf("StringFixed(%v, %d) = %s, want %s", tt.input, tt.precision, got, tt.expected)
		}
	}
}

func TestJSONUnmarshalNull(t *testing.T) {
	var r rat.Rat
	err := r.UnmarshalJSON([]byte(`"null"`))
	if err != nil {
		t.Fatalf("Unmarshal null error: %v", err)
	}
}
