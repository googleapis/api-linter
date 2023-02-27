---
rule:
  aip: 127
  name: [core, '0127', http-template-syntax]
  summary: |
    HTTP patterns should follow the HTTP path template syntax.
permalink: /127/http-template-syntax
redirect_from:
  - /0127/http-template-syntax
---

# HTTP Pattern Variables

This rule enforces that HTTP annotation patterns follow the path template syntax
rules, as mandated in [AIP-127][].

## Details

This rule ensures that `google.api.http` patterns adhere to the following
[syntax rules](https://github.com/googleapis/googleapis/blob/83c3605afb5a39952bf0a0809875d41cf2a558ca/google/api/http.proto#L224).

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc GetBook(GetBookRequest) returns (Book) {
    option (google.api.http) = {
        // Should start with a leading slash.
        get: "v1/{name=shelves/*}"
    };
}
```

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc AddAuthor(AddAuthorRequest) returns (AddAuthorResponse) {
    option (google.api.http) = {
        // Verb should be marked off with the ':' character.
        post: "/v1/{book=publishers/*/books/*}-addAuthor"
        body: "*"
    };
}
```

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    // The triple wildcard ('***') is not a part of the syntax.
    post: "/v1/{parent=publishers/***}"
    body: "book"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc GetBook(GetBookRequest) returns (Book) {
    option (google.api.http) = {
        get: "/v1/{name=shelves/*}"
    };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0127::http-template-syntax=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc GetBook(GetBookRequest) returns (Book) {
    option (google.api.http) = {
        get: "v1/{name=shelves/*}"
    };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-127]: https://aip.dev/127
[aip.dev/not-precedent]: https://aip.dev/not-precedent
