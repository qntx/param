package null

import (
	"bytes"
	"encoding/json"
)

// JSONNullable defines the interface for types that can represent JSON nullability states.
type JSONNullable interface {
	IsNull() bool
	SetNull()
	IsSet() bool
	Reset()
}

// Nullable is a generic type that implements a field with three possible states:
// - field is not set in the request
// - field is explicitly set to `null` in the request
// - field is explicitly set to a valid value in the request
type Nullable[T any] map[bool]T

// Ensure Nullable implements JSONNullable, json.Marshaler, and json.Unmarshaler
var _ JSONNullable = (*Nullable[any])(nil)
var _ json.Marshaler = (*Nullable[any])(nil)
var _ json.Unmarshaler = (*Nullable[any])(nil)

// Zero constructs a Nullable[T] in the unset state, representing a field not provided in a JSON request.
func Zero[T any]() Nullable[T] {
	return make(Nullable[T])
}

// From constructs a Nullable[T] with the given value, representing a field explicitly set in a JSON request.
func From[T any](value T) Nullable[T] {
	return map[bool]T{true: value}
}

// Null constructs a Nullable[T] with an explicit `null`, representing a field set to `null` in a JSON request.
func Null[T any]() Nullable[T] {
	return map[bool]T{false: *new(T)}
}

// Get retrieves the underlying value, if present, and returns an empty value and `false` if not present.
func (t Nullable[T]) Get() (T, bool) {
	var empty T
	if t.IsNull() {
		return empty, false
	}
	if !t.IsSet() {
		return empty, false
	}
	return t[true], true
}

// MustGet retrieves the underlying value, if present, and panics if not present.
func (t Nullable[T]) MustGet() T {
	v, ok := t.Get()
	if !ok {
		panic("value is not set or null")
	}
	return v
}

// Set sets the underlying value to a given value.
func (t *Nullable[T]) Set(value T) {
	*t = map[bool]T{true: value}
}

// IsNull indicates whether the field was sent and had a value of `null`.
func (t Nullable[T]) IsNull() bool {
	_, foundNull := t[false]
	return foundNull
}

// SetNull sets the field to an explicit `null`.
func (t *Nullable[T]) SetNull() {
	*t = map[bool]T{false: *new(T)}
}

// IsSet indicates whether the field was sent (either as null or a value).
func (t Nullable[T]) IsSet() bool {
	return len(t) != 0
}

// Reset clears the field, making it unset.
func (t *Nullable[T]) Reset() {
	*t = map[bool]T{}
}

func (t Nullable[T]) MarshalJSON() ([]byte, error) {
	if t.IsNull() {
		return []byte("null"), nil
	}

	// if field was unspecified, and `omitempty` is set on the field's tags, `json.Marshal` will omit this field

	// otherwise: we have a value, so marshal it
	return json.Marshal(t[true])
}

func (t *Nullable[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		t.SetNull()
		return nil
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	t.Set(v)
	return nil
}
