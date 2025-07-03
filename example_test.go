package null_test

import (
	"encoding/json"
	"fmt"

	"github.com/qntx/null"
)

func ExampleNewNull() {
	p := struct {
		N null.Null[int]
	}{}

	p.N = null.NewNull[int]()

	fmt.Printf("Specified: %v\n", p.N.IsSpecified())
	fmt.Printf("Null: %v\n", p.N.IsNull())
	// Output:
	// Specified: true
	// Null: true
}

func ExampleNewFrom() {
	p := struct {
		N null.Null[int]
	}{}

	p.N = null.NewFrom(123)

	fmt.Println(p.N.Get())
	// Output:
	// 123 <nil>
}

func ExampleNull_marshalRequired() {
	obj := struct {
		ID null.Null[int] `json:"id"`
	}{}

	// when it's not set (by default)
	b, err := json.Marshal(obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Unspecified:")
	fmt.Printf(`JSON: %s`+"\n", b)
	fmt.Println("---")

	// when it's not set (explicitly)
	obj.ID.SetUnspecified()

	b, err = json.Marshal(obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Unspecified:")
	fmt.Printf(`JSON: %s`+"\n", b)
	fmt.Println("---")

	// when it's set explicitly to nil
	obj.ID.SetNull()

	b, err = json.Marshal(obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Null:")
	fmt.Printf(`JSON: %s`+"\n", b)
	fmt.Println("---")

	// when it's set explicitly to the zero value
	var v int
	obj.ID.Set(v)

	b, err = json.Marshal(obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Zero value:")
	fmt.Printf(`JSON: %s`+"\n", b)
	fmt.Println("---")

	// when it's set explicitly to a specific value
	v = 12345
	obj.ID.Set(v)

	b, err = json.Marshal(obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Value:")
	fmt.Printf(`JSON: %s`+"\n", b)
	fmt.Println("---")

	// Output:
	// Unspecified:
	// JSON: {"id":0}
	// ---
	// Unspecified:
	// JSON: {"id":0}
	// ---
	// Null:
	// JSON: {"id":null}
	// ---
	// Zero value:
	// JSON: {"id":0}
	// ---
	// Value:
	// JSON: {"id":12345}
	// ---
}

func ExampleNull_marshalOptional() {
	obj := struct {
		ID null.Null[int] `json:"id,omitempty"`
	}{}

	// when it's not set (by default)
	b, err := json.Marshal(obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Unspecified:")
	fmt.Printf(`JSON: %s`+"\n", b)
	fmt.Println("---")

	// when it's not set (explicitly)
	obj.ID.SetUnspecified()

	b, err = json.Marshal(obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Unspecified:")
	fmt.Printf(`JSON: %s`+"\n", b)
	fmt.Println("---")

	// when it's set explicitly to nil
	obj.ID.SetNull()

	b, err = json.Marshal(obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Null:")
	fmt.Printf(`JSON: %s`+"\n", b)
	fmt.Println("---")

	// when it's set explicitly to the zero value
	var v int
	obj.ID.Set(v)

	b, err = json.Marshal(obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Zero value:")
	fmt.Printf(`JSON: %s`+"\n", b)
	fmt.Println("---")

	// when it's set explicitly to a specific value
	v = 12345
	obj.ID.Set(v)

	b, err = json.Marshal(obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Value:")
	fmt.Printf(`JSON: %s`+"\n", b)
	fmt.Println("---")

	// Output:
	// Unspecified:
	// JSON: {}
	// ---
	// Unspecified:
	// JSON: {}
	// ---
	// Null:
	// JSON: {"id":null}
	// ---
	// Zero value:
	// JSON: {"id":0}
	// ---
	// Value:
	// JSON: {"id":12345}
	// ---
}

func ExampleNull_unmarshalRequired() {
	obj := struct {
		Name null.Null[string] `json:"name"`
	}{}

	// when it's not set
	err := json.Unmarshal([]byte(`
		{
		}
		`), &obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Unspecified:")
	fmt.Printf("obj.Name.IsSpecified(): %v\n", obj.Name.IsSpecified())
	fmt.Printf("obj.Name.IsNull(): %v\n", obj.Name.IsNull())
	fmt.Println("---")

	// when it's set explicitly to nil
	err = json.Unmarshal([]byte(`
		{
		"name": null
		}
		`), &obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Null:")
	fmt.Printf("obj.Name.IsSpecified(): %v\n", obj.Name.IsSpecified())
	fmt.Printf("obj.Name.IsNull(): %v\n", obj.Name.IsNull())
	fmt.Println("---")

	// when it's set explicitly to the zero value
	err = json.Unmarshal([]byte(`
		{
		"name": ""
		}
		`), &obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Zero value:")
	fmt.Printf("obj.Name.IsSpecified(): %v\n", obj.Name.IsSpecified())
	fmt.Printf("obj.Name.IsNull(): %v\n", obj.Name.IsNull())
	val, ok := obj.Name.Get()
	if !ok {
		fmt.Printf("Error: %v\n", "value is not specified or null")
		return
	}
	fmt.Printf("obj.Name.Get(): %#v <nil>\n", val)
	fmt.Printf("obj.Name.MustGet(): %#v\n", obj.Name.MustGet())
	fmt.Println("---")

	// when it's set explicitly to a specific value
	err = json.Unmarshal([]byte(`
		{
		"name": "foo"
		}
		`), &obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Value:")
	fmt.Printf("obj.Name.IsSpecified(): %v\n", obj.Name.IsSpecified())
	fmt.Printf("obj.Name.IsNull(): %v\n", obj.Name.IsNull())
	val, ok = obj.Name.Get()
	if !ok {
		fmt.Printf("Error: %v\n", "value is not specified or null")
		return
	}
	fmt.Printf("obj.Name.Get(): %#v <nil>\n", val)
	fmt.Printf("obj.Name.MustGet(): %#v\n", obj.Name.MustGet())
	fmt.Println("---")

	// Output:
	// Unspecified:
	// obj.Name.IsSpecified(): false
	// obj.Name.IsNull(): false
	// ---
	// Null:
	// obj.Name.IsSpecified(): true
	// obj.Name.IsNull(): true
	// ---
	// Zero value:
	// obj.Name.IsSpecified(): true
	// obj.Name.IsNull(): false
	// obj.Name.Get(): "" <nil>
	// obj.Name.MustGet(): ""
	// ---
	// Value:
	// obj.Name.IsSpecified(): true
	// obj.Name.IsNull(): false
	// obj.Name.Get(): "foo" <nil>
	// obj.Name.MustGet(): "foo"
	// ---
}

func ExampleNull_unmarshalOptional() {
	obj := struct {
		// Note that there is no pointer for null.Null when it's
		Name null.Null[string] `json:"name,omitempty"`
	}{}

	// when it's not set
	err := json.Unmarshal([]byte(`
		{
		}
		`), &obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Unspecified:")
	fmt.Printf("obj.Name.IsSpecified(): %v\n", obj.Name.IsSpecified())
	fmt.Printf("obj.Name.IsNull(): %v\n", obj.Name.IsNull())
	fmt.Println("---")

	// when it's set explicitly to nil
	err = json.Unmarshal([]byte(`
		{
		"name": null
		}
		`), &obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Null:")
	fmt.Printf("obj.Name.IsSpecified(): %v\n", obj.Name.IsSpecified())
	fmt.Printf("obj.Name.IsNull(): %v\n", obj.Name.IsNull())
	fmt.Println("---")

	// when it's set explicitly to the zero value
	err = json.Unmarshal([]byte(`
		{
		"name": ""
		}
		`), &obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Zero value:")
	fmt.Printf("obj.Name.IsSpecified(): %v\n", obj.Name.IsSpecified())
	fmt.Printf("obj.Name.IsNull(): %v\n", obj.Name.IsNull())
	val, ok := obj.Name.Get()
	if !ok {
		fmt.Printf("Error: %v\n", "value is not specified or null")
		return
	}
	fmt.Printf("obj.Name.Get(): %#v <nil>\n", val)
	fmt.Printf("obj.Name.MustGet(): %#v\n", obj.Name.MustGet())
	fmt.Println("---")

	// when it's set explicitly to a specific value
	err = json.Unmarshal([]byte(`
		{
		"name": "foo"
		}
		`), &obj)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Value:")
	fmt.Printf("obj.Name.IsSpecified(): %v\n", obj.Name.IsSpecified())
	fmt.Printf("obj.Name.IsNull(): %v\n", obj.Name.IsNull())
	val, ok = obj.Name.Get()
	if !ok {
		fmt.Printf("Error: %v\n", "value is not specified or null")
		return
	}
	fmt.Printf("obj.Name.Get(): %#v <nil>\n", val)
	fmt.Printf("obj.Name.MustGet(): %#v\n", obj.Name.MustGet())
	fmt.Println("---")

	// Output:
	// Unspecified:
	// obj.Name.IsSpecified(): false
	// obj.Name.IsNull(): false
	// ---
	// Null:
	// obj.Name.IsSpecified(): true
	// obj.Name.IsNull(): true
	// ---
	// Zero value:
	// obj.Name.IsSpecified(): true
	// obj.Name.IsNull(): false
	// obj.Name.Get(): "" <nil>
	// obj.Name.MustGet(): ""
	// ---
	// Value:
	// obj.Name.IsSpecified(): true
	// obj.Name.IsNull(): false
	// obj.Name.Get(): "foo" <nil>
	// obj.Name.MustGet(): "foo"
	// ---
}
