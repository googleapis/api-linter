---
rule:
  aip: 233
  name: [core, '0233', http-method]
  summary: Batch Create methods must use the POST HTTP verb.
---

# Batch Create methods: POST HTTP verb

This rule enforces that all `BatchCreate` RPCs use the `POST` HTTP verb, as
mandated in [AIP-233][].

## Details

This rule looks at any message matching beginning with `BatchCreate`, and
complains if the HTTP verb is anything other than `POST`. It _does_ check
additional bindings if they are present.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc BatchCreateBooks(BatchCreateBooksRequest) returns (BatchCreateBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books:batchCreate" // Should be `post:`.
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc BatchCreateBooks(BatchCreateBooksRequest) returns (BatchCreateBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchCreate"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0233::http-method=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchCreateBooks(BatchCreateBooksRequest) returns (BatchCreateBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books:batchCreate" // Should be `post:`.
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-233]: https://aip.dev/233
[aip.dev/not-precedent]: https://aip.dev/not-precedent
