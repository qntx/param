package null

import (
	"bytes"
	"encoding/json"
)

// Null defines the interface for types that can represent nullability states.
type Null interface {
	IsNull() bool
	SetNull()
	IsSpecified() bool
	SetUnspecified()
}

// Nullable is a generic type, which implements a field that can be one of three states:
//
// - field is not set in the request
// - field is explicitly set to `null` in the request
// - field is explicitly set to a valid value in the request
//
// Nullable is intended to be used with JSON marshalling and unmarshalling.
//
// Internal implementation details:
//
// - map[true]T means a value was provided
// - map[false]T means an explicit null was provided
// - nil or zero map means the field was not provided
//
// If the field is expected to be optional, add the `omitempty` JSON tags. Do NOT use `*Nullable`!
//
// Adapted from https://github.com/golang/go/issues/64515#issuecomment-1841057182
type Nullable[T any] map[bool]T

// Ensure Nullable implements Nullable, json.Marshaler and json.Unmarshaler
var _ Null = (*Nullable[any])(nil)
var _ json.Marshaler = (*Nullable[any])(nil)
var _ json.Unmarshaler = (*Nullable[any])(nil)

// New is a convenience helper to allow constructing a `Nullable` with a given value, for instance to construct a field inside a struct, without introducing an intermediate variable
func New[T any]() Nullable[T] {
	return make(Nullable[T])
}

// NewFrom is a convenience helper to allow constructing a `Nullable` with a given value, for instance to construct a field inside a struct, without introducing an intermediate variable
func NewFrom[T any](t T) Nullable[T] {
	return map[bool]T{true: t}
}

// NewNull is a convenience helper to allow constructing a `Nullable` with an explicit `null`, for instance to construct a field inside a struct, without introducing an intermediate variable
func NewNull[T any]() Nullable[T] {
	return map[bool]T{false: *new(T)}
}

// Get retrieves the underlying value, if present, and returns an empty value and `false` if the value was not present
func (t Nullable[T]) Get() (T, bool) {
	var empty T
	if t.IsNull() {
		return empty, false
	}
	if !t.IsSpecified() {
		return empty, false
	}
	return t[true], true
}

// MustGet retrieves the underlying value, if present, and panics if the value was not present
func (t Nullable[T]) MustGet() T {
	v, ok := t.Get()
	if !ok {
		panic("value is not specified or null")
	}
	return v
}

// Set sets the underlying value to a given value
func (t *Nullable[T]) Set(value T) {
	*t = map[bool]T{true: value}
}

// IsNull indicate whether the field was sent, and had a value of `null`
func (t Nullable[T]) IsNull() bool {
	_, foundNull := t[false]
	return foundNull
}

// SetNull indicate that the field was sent, and had a value of `null`
func (t *Nullable[T]) SetNull() {
	*t = map[bool]T{false: *new(T)}
}

// IsSpecified indicates whether the field was sent
func (t Nullable[T]) IsSpecified() bool {
	return len(t) != 0
}

// SetUnspecified indicate whether the field was sent
func (t *Nullable[T]) SetUnspecified() {
	*t = map[bool]T{}
}

func (t Nullable[T]) MarshalJSON() ([]byte, error) {
	// if field was specified, and `null`, marshal it
	if t.IsNull() {
		return []byte("null"), nil
	}

	// if field was unspecified, and `omitempty` is set on the field's tags, `json.Marshal` will omit this field

	// otherwise: we have a value, so marshal it
	return json.Marshal(t[true])
}

func (t *Nullable[T]) UnmarshalJSON(data []byte) error {
	// if field is unspecified, UnmarshalJSON won't be called

	// if field is specified, and `null`
	if bytes.Equal(data, []byte("null")) {
		t.SetNull()
		return nil
	}
	// otherwise, we have an actual value, so parse it
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	t.Set(v)
	return nil
}
