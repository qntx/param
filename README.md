# null

A Go `Nullable` type for JSON APIs, handling three states: **valid**, **null**, or **unset**. Ideal for `PATCH` requests and optional fields.

## Why?

Go's `encoding/json` struggles to distinguish **absent** fields from **explicit null** in JSON. Existing solutions like `sql.NullString` or `guregu/null` are two-state and fail with `omitempty`. This library uses a generic `map[bool]T` to cleanly handle all three states.

## Features

- **Tri-State**: Differentiates valid, null, and unset.
- **Type-Safe**: Generic support for any type (`[T any]`).
- **Idiomatic**: Simple API, seamless with `omitempty`.
- **Self-Contained**: No custom struct methods needed.

## Installation

```bash
go get github.com/qntx/null
```

## Example

### Struct

```go
type UserPayload struct {
    Name null.Nullable[string] `json:"name,omitempty"`
    Age  null.Nullable[int]    `json:"age,omitempty"`
    Bio  null.Nullable[string] `json:"bio"`
}
```

### Marshaling

```go
payload := UserPayload{
    Name: null.From("Alice"),  // Valid
    Age:  null.Null[int](),    // Null
    Bio:  null.Zero[string](), // Unset
}

data, _ := json.Marshal(payload)
fmt.Println(string(data)) // {"name":"Alice","age":null,"bio":null}

payload2 := UserPayload{
    Name: null.From("Bob"),
    Bio:  null.Null[string](),
}

data, _ = json.Marshal(payload2)
fmt.Println(string(data)) // {"name":"Bob","bio":null}
```

## Performance

`Nullable` is ~3x slower than pointers for marshaling and ~2x for unmarshaling due to its `map` internals. For most APIs, this nanosecond overhead is negligible compared to the clarity and correctness gained.

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
