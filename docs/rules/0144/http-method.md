---
rule:
  aip: 144
  name: [core, '0144', http-method]
  summary: Add/Remove methods must use the POST HTTP verb.
permalink: /144/http-method
redirect_from:
  - /0144/http-method
---

# Add/Remove methods: POST HTTP verb

This rule enforces that all `Add` and `Remove` RPCs use the `POST` HTTP verb, as
mandated in [AIP-144][].

## Details

This rule looks at any RPCs with the name beginning with `Add` or `Remove`, and
complains if the HTTP verb is anything other than `POST`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc AddAuthor(AddAuthorRequest) returns (AddAuthorResponse) {
  option (google.api.http) = {
    patch: "/v1/{book=publishers/*/books/*}:addAuthor" // Should be `post:`.
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
// (-- api-linter: core::0144::http-method=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc AddAuthor(AddAuthorRequest) returns (AddAuthorResponse) {
  option (google.api.http) = {
    patch: "/v1/{book=publishers/*/books/*}:addAuthor" // Should be `post:`.
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-144]: https://aip.dev/144
[aip.dev/not-precedent]: https://aip.dev/not-precedent
