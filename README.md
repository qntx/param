# param

A Go `Opt` type for JSON APIs, handling three states: **valid**, **null**, or **unset**. Ideal for `PATCH` requests and optional fields.

## Why?

Go's `encoding/json` struggles to distinguish **absent** fields from **explicit null** in JSON. Existing solutions like `sql.NullString` or `guregu/null` are two-state and fail with `omitempty`. This library uses a generic `map[bool]T` to cleanly handle all three states.

## Features

- **Tri-State**: Differentiates valid, null, and unset.
- **Type-Safe**: Generic support for any type (`[T any]`).
- **Idiomatic**: Simple API, seamless with `omitempty`.
- **Self-Contained**: No custom struct methods needed.

## Installation

```bash
go get github.com/qntx/param
```

## Example

### Struct

```go
type UserPayload struct {
    Name param.Opt[string] `json:"name,omitempty"`
    Age  param.Opt[int]    `json:"age,omitempty"`
    Bio  param.Opt[string] `json:"bio"`
}
```

### Marshaling

```go
payload := UserPayload{
    Name: param.From("Alice"),  // Valid
    Age:  param.Null[int](),    // Null
    Bio:  param.Zero[string](), // Unset
}

data, _ := json.Marshal(payload)
fmt.Println(string(data)) // {"name":"Alice","age":null,"bio":""}

payload2 := UserPayload{
    Name: param.From("Bob"),
    Bio:  param.Null[string](),
}

data, _ = json.Marshal(payload2)
fmt.Println(string(data)) // {"name":"Bob","bio":null}
```

## Usage Guide

This guide helps you select the appropriate function based on your desired JSON output. The table below maps each field state to its corresponding function and behavior.

| Intent | Function Call | JSON with `omitempty` | JSON without `omitempty` | Common Use Case |
| :--- | :--- | :--- | :--- | :--- |
| **Set a Value** | `param.From("Alice")` | `"field": "Alice"` | `"field": "Alice"` | **Create/Update**: Assign a specific, non-null value. |
| **Set to `null`** | `param.Null[string]()` | `"field": null` | `"field": null` | **Clear**: Explicitly nullify a field's value on the server. |
| **Unset / Omit** | `param.Zero[string]()` | Field is omitted | `"field": ""` (zero-value) | **Partial Update (PATCH)**: Leave a field untouched during an update. |

**Key Points:**

- **`param.From()`**: Use to assign a specific, non-null value to a field.
- **`param.Null()`**: Use to explicitly set a field to `null`.
- **`param.Zero()`**: Ideal for partial updates (`PATCH`). When combined with an `omitempty` tag, the field is excluded from the JSON output, leaving the server-side value unchanged.
- **Important**: Without `omitempty`, `param.Zero()` marshals to the type's zero-value (e.g., `""` for `string`, `0` for `int`).

## Performance

`Opt` is ~3x slower than pointers for marshaling and ~2x for unmarshaling due to its `map` internals. For most APIs, this nanosecond overhead is negligible compared to the clarity and correctness gained.

The following benchmarks were run on a `13th Gen Intel(R) Core(TM) i7-13700H` CPU.

| Benchmark                    | Operations   | Time/Op       | Memory/Op  | Allocations/Op |
| ---------------------------- | ------------ | ------------- | ---------- | -------------- |
| **Marshal (This Library)**   | `3,598,364`  | `336.1 ns/op` | `72 B/op`  | `5 allocs/op`  |
| Marshal (Native Pointers)    | `10,764,397` | `112.0 ns/op` | `48 B/op`  | `2 allocs/op`  |
| **Unmarshal (This Library)** | `1,708,104`  | `707.7 ns/op` | `816 B/op` | `10 allocs/op` |
| Unmarshal (Native Pointers)  | `3,508,210`  | `337.1 ns/op` | `216 B/op` | `4 allocs/op`  |

## License

MIT

## Acknowledgements

- Inspired by [`oapi-codegen/Nullable`](https://github.com/oapi-codegen/Nullable) and the discussion in [Go Issue #64515](https://github.com/golang/go/issues/64515#issuecomment-1841024193).
