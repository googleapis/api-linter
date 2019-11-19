---
rule:
  aip: 234
  name: [core, '0234', plural-method-name]
  summary: Batch Update methods must have plural names.
permalink: /234/plural-method-name
redirect_from:
  - /0234/plural-method-name
---

# Batch Update methods: Plural method name

This rule enforces that all `BatchUpdate` RPCs have a plural resource in the
method name, as mandated in [AIP-234][].

## Details

This rule looks at any method whose name begins with `BatchUpdate`, and complains
if the name of the resource in the method name is singular.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Method name should be `BatchUpdateBooks`
rpc BatchUpdateBook(BatchUpdateBookRequest) returns (BatchUpdateBookResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchUpdate"
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc BatchUpdateBooks(BatchUpdateBooksRequest) returns (BatchUpdateBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchUpdate"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0234::plural-method-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchUpdateBook(BatchUpdateBookRequest) returns (BatchUpdateBookResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchUpdate"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-234]: https://aip.dev/234
[aip.dev/not-precedent]: https://aip.dev/not-precedent
