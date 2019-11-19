---
rule:
  aip: 234
  name: [core, '0234', http-method]
  summary: Batch Update methods must use the POST HTTP verb.
permalink: /234/http-method
redirect_from:
  - /0234/http-method
---

# Batch Update methods: POST HTTP verb

This rule enforces that all `BatchUpdate` RPCs use the `POST` HTTP verb, as
mandated in [AIP-234][].

## Details

This rule looks at any RPCs with the name beginning with `BatchUpdate`, and
complains if the HTTP verb is anything other than `POST`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc BatchUpdateBooks(BatchUpdateBooksRequest) returns (BatchUpdateBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books:batchUpdate" // Should be `post:`.
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
// (-- api-linter: core::0234::http-method=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchUpdateBooks(BatchUpdateBooksRequest) returns (BatchUpdateBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books:batchUpdate" // Should be `post:`.
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-234]: https://aip.dev/234
[aip.dev/not-precedent]: https://aip.dev/not-precedent
