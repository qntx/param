package null

import (
	"bytes"
	"encoding/json"
)

// State represents the possible states of a Null value.
type State uint8

// State constants for Null type.
const (
	StateUnset State = iota // 0: Value is not set
	StateNull               // 1: Value is explicitly null
	StateValid              // 2: Value is valid
)

// Nullable defines the interface for types that can represent nullability states.
type Nullable interface {
	IsNull() bool
	SetNull()
	IsSpecified() bool
	IsZero() bool
	SetUnspecified()
}

// Null is a generic type that represents a field with three possible states:
// - Not set (unset)
// - Explicitly set to null
// - Explicitly set to a valid value
//
// It implements json.Marshaler and json.Unmarshaler for JSON handling.
// Use the `omitempty` JSON tag for optional fields.
type Null[T any] struct {
	value T     // The underlying value
	state State // Current state
}

// Ensure Null implements Nullable, json.Marshaler and json.Unmarshaler
var _ Nullable = (*Null[any])(nil)
var _ json.Marshaler = (*Null[any])(nil)
var _ json.Unmarshaler = (*Null[any])(nil)

// New creates a Null in the unset state.
func New[T any]() Null[T] {
	return Null[T]{}
}

// NewFrom creates a Null with a valid value.
func NewFrom[T any](t T) Null[T] {
	return Null[T]{
		value: t,
		state: StateValid,
	}
}

// NewNull creates a Null in the null state.
func NewNull[T any]() Null[T] {
	return Null[T]{
		state: StateNull,
	}
}

// Get retrieves the value and a boolean indicating if it's valid.
func (n Null[T]) Get() (T, bool) {
	if n.state != StateValid {
		var zero T
		return zero, false
	}
	return n.value, true
}

// MustGet retrieves the value or panics if not valid.
func (n Null[T]) MustGet() T {
	if v, ok := n.Get(); ok {
		return v
	}
	panic("value is not specified")
}

// Set sets the value and marks it as valid.
func (n *Null[T]) Set(value T) {
	n.value = value
	n.state = StateValid
}

// IsNull checks if the value is explicitly null.
func (n Null[T]) IsNull() bool {
	return n.state == StateNull
}

// SetNull sets the value to null.
func (n *Null[T]) SetNull() {
	n.value = *new(T)
	n.state = StateNull
}

// IsSpecified checks if the value is set (null or valid).
func (n Null[T]) IsSpecified() bool {
	return n.state != StateUnset
}

// IsZero returns true if the field is unset, supporting omitempty.
func (n Null[T]) IsZero() bool {
	return n.state == StateUnset
}

// SetUnspecified sets the value to unset.
func (n *Null[T]) SetUnspecified() {
	n.value = *new(T)
	n.state = StateUnset
}

// MarshalJSON implements json.Marshaler.
func (n Null[T]) MarshalJSON() ([]byte, error) {
	if n.state == StateNull {
		return []byte("null"), nil
	}
	// if field was unspecified, and `omitempty` is not set on the field's tags, `json.Marshal` will include this field
	// Unset fields with `omitempty` tag will be omitted by json.Marshal
	// otherwise: we have a value, so marshal it
	return json.Marshal(n.value)
}

// UnmarshalJSON implements json.Unmarshaler.
func (n *Null[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		n.SetNull()
		return nil
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	n.Set(v)

	return nil
}
