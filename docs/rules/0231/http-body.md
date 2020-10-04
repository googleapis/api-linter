---
rule:
  aip: 231
  name: [core, '0231', http-body]
  summary: Batch Get methods must not have an HTTP body.
permalink: /231/http-body
redirect_from:
  - /0231/http-body
---

# Batch Get methods: No HTTP body

This rule enforces that all `BatchGet` RPCs omit the HTTP `body`, as mandated in
[AIP-231][].

## Details

This rule looks at any method beginning with `BatchGet`, and
complains if the HTTP `body` field is set.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc BatchGetBooks(BatchGetBooksRequest) returns (BatchGetBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchGet"
    body: "*"  // This should be absent.
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
// (-- api-linter: core::0231::http-body=disabled
//     api-linter: core::0231::http-method=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchGetBooks(BatchGetBooksRequest) returns (BatchGetBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchGet"
    body: "*"  // This should be absent.
  };
}
```

**Important:** HTTP `GET` requests are unable to have an HTTP body, due to the
nature of the protocol. The only valid way to include a body is to also use a
different HTTP method (as depicted above).

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-231]: https://aip.dev/231
[aip.dev/not-precedent]: https://aip.dev/not-precedent
