---
rule:
  aip: 231
  name: [core, '0231', http-uri-suffix]
  summary: Batch Get methods must have the correct URI suffix
permalink: /231/http-uri-suffix
redirect_from:
  - /0231/http-uri-suffix
---

# Batch Get methods: URI suffix

This rule enforces that ` BatchGet` methods include the `:batchGet` suffix in
the REST URI, as mandated in [AIP-231][].

## Details

This rule looks at any method whose name starts with `BatchGet`, and complains
if the HTTP URI does not end with `:batchGet`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc BatchGetBooks(BatchGetBooksRequest) returns (BatchGetBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books:batch"  // Should end with `:batchGet`.
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
// (-- api-linter: core::0231::http-uri-suffix=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchGetBooks(BatchGetBooksRequest) returns (BatchGetBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books:batch"  // Should end with `:batchGet`.
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-231]: https://aip.dev/231
[aip.dev/not-precedent]: https://aip.dev/not-precedent
