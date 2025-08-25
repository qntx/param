package param_test

import (
	"encoding/json"
	"testing"

	"github.com/qntx/param"
)

// blackhole is used to prevent the compiler from optimizing away benchmark results.
var blackhole any

// OptPayload uses the custom null.Opt[T] type.
type OptPayload struct {
	ID    param.Opt[int]    `json:"id"`
	Name  param.Opt[string] `json:"name"`
	Email param.Opt[string] `json:"email,omitempty"` // For omitempty test
}

// PointerPayload uses the standard Go pointer types for nullability.
type PointerPayload struct {
	ID    *int    `json:"id"`
	Name  *string `json:"name"`
	Email *string `json:"email,omitempty"`
}

var (
	// A representative JSON input for unmarshaling benchmarks.
	// It includes a valid value, an explicit null, and a missing optional field.
	jsonInput = []byte(`{"id":12345,"name":null}`)

	// A sample ID value for struct construction.
	idValue = 12345

	// Pre-constructed structs for marshaling benchmarks.

	// Data using the custom Opt type.
	OptData = OptPayload{
		ID:    param.From(idValue),
		Name:  param.Null[string](),
		Email: nil, // Unset, will be omitted.
	}

	// Data using standard pointer types.
	pointerData = PointerPayload{
		ID:    &idValue,
		Name:  nil, // In pointers, a nil value marshals to `null`.
		Email: nil, // A nil value with `omitempty` is omitted.
	}
)

// BenchmarkMarshal_WithOpt tests the performance of marshaling a struct that uses null.Opt[T].
func BenchmarkMarshal_WithOpt(b *testing.B) {
	b.ReportAllocs()
	var r []byte
	var err error

	for i := 0; i < b.N; i++ {
		r, err = json.Marshal(OptData)
	}

	blackhole = r
	if err != nil {
		b.Fatal(err)
	}
}

// BenchmarkMarshal_WithPointers tests the performance of marshaling a struct using standard pointers.
func BenchmarkMarshal_WithPointers(b *testing.B) {
	b.ReportAllocs()
	var r []byte
	var err error

	for i := 0; i < b.N; i++ {
		r, err = json.Marshal(pointerData)
	}

	blackhole = r
	if err != nil {
		b.Fatal(err)
	}
}

// BenchmarkUnmarshal_WithOpt tests the performance of unmarshaling into a struct that uses null.Opt[T].
func BenchmarkUnmarshal_WithOpt(b *testing.B) {
	b.ReportAllocs()
	var p OptPayload
	var err error

	for i := 0; i < b.N; i++ {
		err = json.Unmarshal(jsonInput, &p)
	}

	blackhole = p
	if err != nil {
		b.Fatal(err)
	}
}

// BenchmarkUnmarshal_WithPointers tests the performance of unmarshaling into a struct using standard pointers.
func BenchmarkUnmarshal_WithPointers(b *testing.B) {
	b.ReportAllocs()
	var p PointerPayload
	var err error

	for i := 0; i < b.N; i++ {
		err = json.Unmarshal(jsonInput, &p)
	}

	blackhole = p
	if err != nil {
		b.Fatal(err)
	}
}
