# Null

> An implementation of a `Null` type for JSON bodies, indicating whether the field is absent, set to null, or set to a value

Unlike other known implementations, this makes it possible to both marshal and unmarshal the value, as well as represent all three states:

- the field is _not set_
- the field is _explicitly set to null_
- the field is _explicitly set to a given value_

And can be embedded in structs, for instance with the following definition:

```go
obj := struct {
    // RequiredID is a required, Null field
    RequiredID     null.Null[int]     `json:"id"`
    // OptionalString is an optional, Null field
    // NOTE that no pointer is required, only `omitempty`
    OptionalString null.Null[string] `json:"optionalString,omitempty"`
}{}
```

## License

MIT

## Acknowledgments

- [Nullable](https://github.com/oapi-codegen/Nullable)  
