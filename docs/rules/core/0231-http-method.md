---
rule:
  aip: 231
  name: [core, '0231', http-method]
  summary: Batch Get methods must use the GET HTTP verb.
---

# Get methods: GET HTTP verb

This rule enforces that all `BatchGet` RPCs use the `GET` HTTP verb, as
mandated in [AIP-231][].

## Details

This rule looks at any message matching beginning with `BatchGet`, and
complains if the HTTP verb is anything other than `GET`. It _does_ check
additional bindings if they are present.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc BatchGetBooks(BatchGetBooksRequest) returns (BatchGetBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*}/books:batchGet"  // Should be `get:`.
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc BatchGetBooks(BatchGetBooksRequest) returns (BatchGetBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*}/books:batchGet"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0231::http-method=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchGetBooks(BatchGetBooksRequest) returns (BatchGetBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*}/books:batchGet"  // Should be `get:`.
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-231]: https://aip.dev/231
[aip.dev/not-precedent]: https://aip.dev/not-precedent
