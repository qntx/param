package null_test

import (
	"encoding/json"
	"testing"

	"github.com/qntx/null" // Assuming this is your local package path
	"github.com/stretchr/testify/require"
)

// Obj is the struct used for testing.
// The `omitempty` tag is important for optional fields.
type Obj struct {
	Foo null.Null[string] `json:"foo,omitempty"`
}

// TestNull_JSONUnmarshal tests all scenarios of deserializing from JSON.
func TestNull_JSONUnmarshal(t *testing.T) {
	t.Run("Value Present", func(t *testing.T) {
		data := `{"foo":"bar"}`
		expectedObj := Obj{Foo: null.NewFrom("bar")}

		// Deserialize from JSON.
		myObj := parse(data, t)
		require.Equal(t, expectedObj, myObj)

		// Verify the state.
		require.True(t, myObj.Foo.IsSpecified())
		require.False(t, myObj.Foo.IsNull())

		// Verify value retrieval.
		value, ok := myObj.Foo.Get()
		require.True(t, ok)
		require.Equal(t, "bar", value)
		require.Equal(t, "bar", myObj.Foo.MustGet())

		// Serialize back to JSON.
		require.JSONEq(t, data, serialize(myObj, t))
	})

	t.Run("Value Absent (omitempty)", func(t *testing.T) {
		data := `{}`
		// When the corresponding field is absent in the JSON, we get a zero value in the Unset state.
		expectedObj := Obj{Foo: null.New[string]()}

		// Deserialize from JSON.
		myObj := parse(data, t)
		require.Equal(t, expectedObj, myObj)

		// Verify the state.
		require.False(t, myObj.Foo.IsSpecified())
		require.False(t, myObj.Foo.IsNull())
		require.True(t, myObj.Foo.IsZero()) // This ensures `omitempty` works correctly.

		// Verify value retrieval.
		_, ok := myObj.Foo.Get()
		require.False(t, ok)

		// Serialize back to JSON.
		require.JSONEq(t, data, serialize(myObj, t))
	})

	t.Run("Value is Null", func(t *testing.T) {
		data := `{"foo":null}`
		expectedObj := Obj{Foo: null.NewNull[string]()}

		// Deserialize from JSON.
		myObj := parse(data, t)
		require.Equal(t, expectedObj, myObj)

		// Verify the state.
		require.True(t, myObj.Foo.IsSpecified())
		require.True(t, myObj.Foo.IsNull())
		require.False(t, myObj.Foo.IsZero())

		// Verify value retrieval.
		_, ok := myObj.Foo.Get()
		require.False(t, ok)
		require.Panics(t, func() { myObj.Foo.MustGet() })

		// Serialize back to JSON.
		require.JSONEq(t, data, serialize(myObj, t))
	})
}

// TestNull_ProgrammaticCreation tests scenarios of creating objects programmatically in Go.
func TestNull_ProgrammaticCreation(t *testing.T) {
	t.Run("With Value", func(t *testing.T) {
		// Using the constructor.
		myObj1 := Obj{Foo: null.NewFrom("bar")}
		require.JSONEq(t, `{"foo":"bar"}`, serialize(myObj1, t))

		// Using the Set method.
		var myObj2 Obj
		myObj2.Foo.Set("bar")
		require.JSONEq(t, `{"foo":"bar"}`, serialize(myObj2, t))
	})

	t.Run("Unspecified", func(t *testing.T) {
		// The default zero value is the Unset state.
		var myObj Obj
		require.True(t, myObj.Foo.IsZero())
		require.JSONEq(t, `{}`, serialize(myObj, t))

		// Explicitly set to Unset.
		myObj.Foo.Set("some value") // Start with a value.
		myObj.Foo.SetUnspecified()  // Then set it back to Unset.
		require.True(t, myObj.Foo.IsZero())
		require.JSONEq(t, `{}`, serialize(myObj, t))
	})

	t.Run("With Null", func(t *testing.T) {
		// Using the constructor.
		myObj1 := Obj{Foo: null.NewNull[string]()}
		require.JSONEq(t, `{"foo":null}`, serialize(myObj1, t))

		// Using the SetNull method.
		var myObj2 Obj
		myObj2.Foo.SetNull()
		require.JSONEq(t, `{"foo":null}`, serialize(myObj2, t))
	})
}

// parse is a helper function to deserialize an object from a string.
func parse(data string, t *testing.T) Obj {
	var myObj Obj
	err := json.Unmarshal([]byte(data), &myObj)
	require.NoError(t, err)
	return myObj
}

// serialize is a helper function to serialize an object to a string.
func serialize(o Obj, t *testing.T) string {
	data, err := json.Marshal(o)
	require.NoError(t, err)
	return string(data)
}
