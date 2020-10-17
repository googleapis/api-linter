---
rule:
  aip: 144
  name: [core, '0144', request-message-name]
  summary: Add/Remove methods must have standardized request message names.
permalink: /144/request-message-name
redirect_from:
  - /0144/request-message-name
---

# Add/Remove methods: Request message

This rule enforces that all `Add` and `Remove` RPCs have a request message name
of `Add*Request` or `Remove*Request`, as mandated in [AIP-144][].

## Details

This rule looks at any message beginning with `Add` or `Remove`, and complains
if the name of the corresponding input message does not match the name of the
RPC with the suffix `Request` appended.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `AddAuthorRequest`.
rpc AddAuthor(AppendAuthorRequest) returns (AddAuthorResponse) {
  option (google.api.http) = {
    post: "/v1/{book=publishers/*/books/*}:addAuthor"
    body: "*"
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
// (-- api-linter: core::0144::request-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc AddAuthor(AppendAuthorRequest) returns (AddAuthorResponse) {
  option (google.api.http) = {
    post: "/v1/{book=publishers/*/books/*}:addAuthor"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-144]: https://aip.dev/144
[aip.dev/not-precedent]: https://aip.dev/not-precedent
