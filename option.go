package param

import (
	"bytes"
	"encoding/json"
)

// JSONOpt defines the interface for types that can represent JSON nullability states.
type JSONOpt interface {
	IsNull() bool
	SetNull()
	IsSet() bool
	Reset()
}

// Opt is a generic type that implements a field with three possible states:
// - field is not set in the request
// - field is explicitly set to `null` in the request
// - field is explicitly set to a valid value in the request
type Opt[T any] map[bool]T

// Ensure Opt implements JSONOpt, json.Marshaler, and json.Unmarshaler
var _ JSONOpt = (*Opt[any])(nil)
var _ json.Marshaler = (*Opt[any])(nil)
var _ json.Unmarshaler = (*Opt[any])(nil)

// Zero constructs a Opt[T] in the unset state, representing a field not provided in a JSON request.
func Zero[T any]() Opt[T] {
	return make(Opt[T])
}

// From constructs a Opt[T] with the given value, representing a field explicitly set in a JSON request.
func From[T any](value T) Opt[T] {
	return map[bool]T{true: value}
}

// Null constructs a Opt[T] with an explicit `null`, representing a field set to `null` in a JSON request.
func Null[T any]() Opt[T] {
	return map[bool]T{false: *new(T)}
}

// Get retrieves the underlying value, if present, and returns an empty value and `false` if not present.
func (t Opt[T]) Get() (T, bool) {
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
func (t Opt[T]) MustGet() T {
	v, ok := t.Get()
	if !ok {
		panic("value is not set or null")
	}
	return v
}

// Set sets the underlying value to a given value.
func (t *Opt[T]) Set(value T) {
	*t = map[bool]T{true: value}
}

// IsNull indicates whether the field was sent and had a value of `null`.
func (t Opt[T]) IsNull() bool {
	_, foundNull := t[false]
	return foundNull
}

// SetNull sets the field to an explicit `null`.
func (t *Opt[T]) SetNull() {
	*t = map[bool]T{false: *new(T)}
}

// IsSet indicates whether the field was sent (either as null or a value).
func (t Opt[T]) IsSet() bool {
	return len(t) != 0
}

// Reset clears the field, making it unset.
func (t *Opt[T]) Reset() {
	*t = map[bool]T{}
}

func (t Opt[T]) MarshalJSON() ([]byte, error) {
	if t.IsNull() {
		return []byte("null"), nil
	}

	// if field was unspecified, and `omitempty` is set on the field's tags, `json.Marshal` will omit this field

	// otherwise: we have a value, so marshal it
	return json.Marshal(t[true])
}

func (t *Opt[T]) UnmarshalJSON(data []byte) error {
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
