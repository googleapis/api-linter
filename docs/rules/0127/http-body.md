---
rule:
  aip: 127
  name: [core, '0127', http-body]
  summary: `GET` and `DELETE` methods must not have an HTTP `body`.
permalink: /127/http-body
redirect_from:
  - /0127/http-body
  - /131/http-body
  - /0131/http-body
  - /132/http-body
  - /0132/http-body
  - /135/http-body
  - /0135/http-body
  - /162/delete-revision-http-body
  - /0162/delete-revision-http-body
  - /162/list-revisions-http-body
  - /0162/list-revisions-http-body
  - /231/http-body
  - /0231/http-body
---

# HTTP URI case

This rule enforces that `GET` and `DELETE` methods omit the HTTP `body`, as
mandated in [AIP-127](http://aip.dev/127).

## Details

This rule looks at all HTTP rules that use `GET` or `DELETE` as the verb,
and complains if the `body` field is set.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc GetBook(GetBookRequest) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}"
    body: "*"  // This should be omitted.
  };
}

// Incorrect.
rpc DeleteBook(DeleteBookRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}"
    body: "*"  // This should be omitted.
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc GetBook(GetBookRequest) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}"
  };
}

// Correct.
rpc DeleteBook(DeleteBookRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0127::http-body=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc GetBook(GetBookRequest) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-127]: https://aip.dev/127
[aip.dev/not-precedent]: https://aip.dev/not-precedent
