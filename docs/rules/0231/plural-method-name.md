---
rule:
  aip: 231
  name: [core, '0231', plural-method-name]
  summary: Batch Get methods must have plural names.
permalink: /231/plural-method-name
redirect_from:
  - /0231/plural-method-name
---

# Batch Get methods: Plural method name

This rule enforces that all `BatchGet` RPCs have a plural resource in the
method name, as mandated in [AIP-231][].

## Details

This rule looks at any method whose name begins with `BatchGet`, and complains
if the name of the resource in the method name is singular.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Method name should be `BatchGetBooks`
rpc BatchGetBook(BatchGetBookRequest) returns (BatchGetBookResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books:batchGet"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc BatchGetBooks(BatchGetBooksRequest) returns (BatchGetBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books:batchGet"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0231::plural-method-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchGetBook(BatchGetBookRequest) returns (BatchGetBookResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books:batchGet"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-231]: https://aip.dev/231
[aip.dev/not-precedent]: https://aip.dev/not-precedent
