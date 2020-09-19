---
rule:
  aip: 235
  name: [core, '0235', plural-method-name]
  summary: Batch Delete methods must have plural names.
permalink: /235/plural-method-name
redirect_from:
  - /0235/plural-method-name
---

# Batch Delete methods: Plural method name

This rule enforces that all `BatchDelete` RPCs have a plural resource in the
method name, as mandated in [AIP-235][].

## Details

This rule looks at any method whose name begins with `BatchDelete`, and complains
if the name of the resource in the method name is singular.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Method name should be `BatchDeleteBooks`
rpc BatchDeleteBook(BatchDeleteBookRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchDelete"
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc BatchDeleteBooks(BatchDeleteBooksRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchDelete"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0235::plural-method-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchDeleteBook(BatchDeleteBookRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchDelete"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-235]: https://aip.dev/235
[aip.dev/not-precedent]: https://aip.dev/not-precedent
