package param_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/qntx/param"
)

// TestConstructors verifies the behavior of the New, From, and Null functions.
func TestConstructors(t *testing.T) {
	t.Run("From should create a valid Opt", func(t *testing.T) {
		n := param.From("hello")

		if n.IsNull() {
			t.Error("Expected IsNull to be false")
		}
		if !n.IsSet() {
			t.Error("Expected IsSet to be true")
		}
		if val, ok := n.Get(); !ok || val != "hello" {
			t.Errorf(`Get() got (%q, %v), want ("hello", true)`, val, ok)
		}
	})

	t.Run("Null should create a null Opt", func(t *testing.T) {
		n := param.Null[int]()

		if !n.IsNull() {
			t.Error("Expected IsNull to be true")
		}
		if !n.IsSet() {
			t.Error("Expected IsSet to be true")
		}
		if _, ok := n.Get(); ok {
			t.Error("Get() on a null value should return false")
		}
	})

	t.Run("New should create an empty (unspecified) Opt", func(t *testing.T) {
		n := param.Zero[any]()

		if n.IsNull() {
			t.Error("Expected IsNull to be false")
		}
		if n.IsSet() {
			t.Error("Expected IsSet to be false for a new empty Opt")
		}
	})
}

// TestStateChecks validates the IsSet and IsNull methods across all states.
func TestStateChecks(t *testing.T) {
	testCases := []struct {
		name        string
		n           param.Opt[any]
		isSpecified bool
		isNull      bool
	}{
		{
			name:        "Unset (nil map)",
			n:           nil,
			isSpecified: false,
			isNull:      false,
		},
		{
			name:        "Unset (empty map)",
			n:           param.Zero[any](),
			isSpecified: false,
			isNull:      false,
		},
		{
			name:        "Null",
			n:           param.Null[any](),
			isSpecified: true,
			isNull:      true,
		},
		{
			name:        "Valid",
			n:           param.From[any]("value"),
			isSpecified: true,
			isNull:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.n.IsSet(); got != tc.isSpecified {
				t.Errorf("IsSet() got %v, want %v", got, tc.isSpecified)
			}
			if got := tc.n.IsNull(); got != tc.isNull {
				t.Errorf("IsNull() got %v, want %v", got, tc.isNull)
			}
		})
	}
}

// TestGetters validates the Get and MustGet methods.
func TestGetters(t *testing.T) {
	t.Run("Get on valid value", func(t *testing.T) {
		n := param.From(42)
		val, ok := n.Get()
		if !ok {
			t.Fatal("Get() returned ok=false for a valid value")
		}
		if val != 42 {
			t.Errorf("Get() value got %d, want 42", val)
		}
	})

	t.Run("Get on null value", func(t *testing.T) {
		n := param.Null[int]()
		_, ok := n.Get()
		if ok {
			t.Error("Get() returned ok=true for a null value")
		}
	})

	t.Run("Get on unset value", func(t *testing.T) {
		var n param.Opt[int]
		_, ok := n.Get()
		if ok {
			t.Error("Get() returned ok=true for an unset value")
		}
	})

	t.Run("MustGet panics", func(t *testing.T) {
		testCases := map[string]param.Opt[string]{
			"for null value":  param.Null[string](),
			"for unset value": nil,
		}

		for name, n := range testCases {
			t.Run(name, func(t *testing.T) {
				defer func() {
					if r := recover(); r == nil {
						t.Error("MustGet should have panicked but did not")
					}
				}()
				_ = n.MustGet() // This line should panic
			})
		}
	})
}

// TestSetters validates the Set, SetNull, and SetUnspecified methods.
func TestSetters(t *testing.T) {
	t.Run("Set should make value valid", func(t *testing.T) {
		var n param.Opt[string]
		n.Set("new value")
		if !n.IsSet() || n.IsNull() {
			t.Error("Set() failed to make value specified and valid")
		}
		if val, _ := n.Get(); val != "new value" {
			t.Errorf("Get() after Set() got %q, want 'new value'", val)
		}
	})

	t.Run("SetNull should make value null", func(t *testing.T) {
		n := param.From("initial value")
		n.SetNull()
		if !n.IsSet() || !n.IsNull() {
			t.Error("SetNull() failed to make value specified and null")
		}
	})

	t.Run("SetUnspecified should make value unspecified", func(t *testing.T) {
		n := param.From("initial value")
		n.Reset()
		if n.IsSet() {
			t.Error("SetUnspecified() failed to make value unspecified")
		}
	})
}

// TestJSONMarshaling validates the JSON serialization logic.
func TestJSONMarshaling(t *testing.T) {
	type Payload struct {
		Required param.Opt[string] `json:"required"`
		Optional param.Opt[int]    `json:"optional,omitempty"`
		Always   param.Opt[bool]   `json:"always"`
	}

	testCases := []struct {
		name  string
		input Payload
		want  string
	}{
		{
			name: "All fields valid",
			input: Payload{
				Required: param.From("hello"),
				Optional: param.From(123),
				Always:   param.From(true),
			},
			want: `{"required":"hello","optional":123,"always":true}`,
		},
		{
			name: "Optional field is unset and omitted",
			input: Payload{
				Required: param.From("world"),
				Optional: nil, // Unset, so it should be omitted
				Always:   param.Null[bool](),
			},
			want: `{"required":"world","always":null}`,
		},
		{
			name: "Required field is unset (marshals to null)",
			input: Payload{
				Required: nil, // Unset, but without omitempty
				Optional: param.From(42),
				Always:   param.From(false),
			},
			want: `{"required":"","optional":42,"always":false}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.input)
			if err != nil {
				t.Fatalf("json.Marshal() returned an unexpected error: %v", err)
			}
			assertJSONEquals(t, got, []byte(tc.want))
		})
	}
}

// TestJSONUnmarshaling validates the JSON deserialization logic.
func TestJSONUnmarshaling(t *testing.T) {
	type Payload struct {
		Required param.Opt[string] `json:"required"`
		Optional param.Opt[int]    `json:"optional"`
	}

	t.Run("All fields present and valid", func(t *testing.T) {
		input := `{"required":"hello","optional":123}`
		var p Payload
		if err := json.Unmarshal([]byte(input), &p); err != nil {
			t.Fatalf("json.Unmarshal() failed: %v", err)
		}
		if p.Required.MustGet() != "hello" || p.Optional.MustGet() != 123 {
			t.Error("Did not unmarshal valid values correctly")
		}
	})

	t.Run("Optional field is missing", func(t *testing.T) {
		input := `{"required":"world"}`
		var p Payload
		if err := json.Unmarshal([]byte(input), &p); err != nil {
			t.Fatalf("json.Unmarshal() failed: %v", err)
		}
		if !p.Optional.IsSet() == false {
			t.Error("A missing field should be Unspecified")
		}
		if p.Required.MustGet() != "world" {
			t.Error("Did not unmarshal valid required value correctly")
		}
	})

	t.Run("Fields are explicitly null", func(t *testing.T) {
		input := `{"required":null,"optional":null}`
		var p Payload
		if err := json.Unmarshal([]byte(input), &p); err != nil {
			t.Fatalf("json.Unmarshal() failed: %v", err)
		}
		if !p.Required.IsNull() || !p.Optional.IsNull() {
			t.Error("Fields should be Null after unmarshaling JSON null")
		}
	})

	t.Run("Type mismatch error", func(t *testing.T) {
		input := `{"required":123}` // required is a string
		var p Payload
		if err := json.Unmarshal([]byte(input), &p); err == nil {
			t.Fatal("Expected a type mismatch error but got nil")
		}
	})
}

// assertJSONEquals is a helper to compare two JSON byte slices by comparing their
// map representations, which ignores key ordering differences.
func assertJSONEquals(t *testing.T, got, want []byte) {
	t.Helper()

	var gotMap, wantMap map[string]interface{}

	if err := json.Unmarshal(got, &gotMap); err != nil {
		t.Fatalf("Failed to unmarshal 'got' JSON: %v\nJSON: %s", err, got)
	}
	if err := json.Unmarshal(want, &wantMap); err != nil {
		t.Fatalf("Failed to unmarshal 'want' JSON: %v\nJSON: %s", err, want)
	}

	if !reflect.DeepEqual(gotMap, wantMap) {
		// For better error reporting, re-marshal with indentation
		gotFormatted, _ := json.MarshalIndent(gotMap, "", "  ")
		wantFormatted, _ := json.MarshalIndent(wantMap, "", "  ")
		t.Errorf("JSON mismatch:\n--- GOT:\n%s\n--- WANT:\n%s\n", gotFormatted, wantFormatted)
	}
}
