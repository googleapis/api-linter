---
rule:
  aip: 144
  name: [core, '0144', http-body]
  summary: Add/Remove methods should use `*` as the HTTP body.
permalink: /144/http-body
redirect_from:
  - /0144/http-body
---

# Add/Remove methods: HTTP body

This rule enforces that all `Add` and `Remove` RPCs use `*` as the HTTP `body`, as
mandated in [AIP-144][].

## Details

This rule looks at any RPC methods beginning with `Add` or `Remove`, and
complains if the HTTP `body` field is anything other than `*`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc AddAuthor(AddAuthorRequest) returns (AddAuthorResponse) {
  option (google.api.http) = {
    post: "/v1/{book=publishers/*/books/*}:addAuthor"
    body: ""  // The http body should be "*".
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc AddAuthor(AddAuthorRequest) returns (AddAuthorResponse) {
  option (google.api.http) = {
    post: "/v1/{book=publishers/*/books/*}:addAuthor"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0144::http-body=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc AddAuthor(AddAuthorRequest) returns (AddAuthorResponse) {
  option (google.api.http) = {
    post: "/v1/{book=publishers/*/books/*}:addAuthor"
    body: ""
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-144]: https://aip.dev/144
[aip.dev/not-precedent]: https://aip.dev/not-precedent
