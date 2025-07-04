package null_test

import (
	"encoding/json"
	"testing"

	"github.com/qntx/null"
)

// blackhole is used to prevent the compiler from optimizing away benchmark results.
var blackhole any

// NullablePayload uses the custom null.Nullable[T] type.
type NullablePayload struct {
	ID    null.Nullable[int]    `json:"id"`
	Name  null.Nullable[string] `json:"name"`
	Email null.Nullable[string] `json:"email,omitempty"` // For omitempty test
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

	// Data using the custom Nullable type.
	nullableData = NullablePayload{
		ID:    null.From(idValue),
		Name:  null.Null[string](),
		Email: nil, // Unset, will be omitted.
	}

	// Data using standard pointer types.
	pointerData = PointerPayload{
		ID:    &idValue,
		Name:  nil, // In pointers, a nil value marshals to `null`.
		Email: nil, // A nil value with `omitempty` is omitted.
	}
)

// BenchmarkMarshal_WithNullable tests the performance of marshaling a struct that uses null.Nullable[T].
func BenchmarkMarshal_WithNullable(b *testing.B) {
	b.ReportAllocs()
	var r []byte
	var err error

	for i := 0; i < b.N; i++ {
		r, err = json.Marshal(nullableData)
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

// BenchmarkUnmarshal_WithNullable tests the performance of unmarshaling into a struct that uses null.Nullable[T].
func BenchmarkUnmarshal_WithNullable(b *testing.B) {
	b.ReportAllocs()
	var p NullablePayload
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
