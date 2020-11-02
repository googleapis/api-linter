---
rule:
  aip: 132
  name: [core, '0132', http-body]
  summary: List methods must not have an HTTP body.
permalink: /132/http-body
redirect_from:
  - /0132/http-body
---

# List methods: No HTTP body

This rule enforces that all `List` RPCs omit the HTTP `body`, as mandated in
[AIP-132][].

## Details

This rule looks at any message beginning with `List`, and complains if
the HTTP `body` field is set.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books"
    body: "*"  // This should be absent.
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0132::http-body=disabled
//     api-linter: core::0132::http-method=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books"
    body: "*"
  };
}
```

**Important:** HTTP `GET` requests are unable to have an HTTP body, due to the
nature of the protocol. The only valid way to include a body is to also use a
different HTTP method (as depicted above).

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-132]: https://aip.dev/132
[aip.dev/not-precedent]: https://aip.dev/not-precedent
