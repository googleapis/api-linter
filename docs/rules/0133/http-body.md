---
rule:
  aip: 133
  name: [core, '0133', http-body]
  summary: Create methods must have the HTTP body set to the resource.
permalink: /133/http-body
redirect_from:
  - /0133/http-body
---

# Create methods: HTTP body

This rule enforces that all `Create` RPCs set the HTTP `body` to the resource,
as mandated in [AIP-133][].

## Details

This rule looks at any message matching beginning with `Create`, and complains
if the HTTP `body` field is not set to the resource being created.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books"
    body: "*"  // This should be "book".
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books"
    body: "book"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0133::http-body=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-133]: https://aip.dev/133
[aip.dev/not-precedent]: https://aip.dev/not-precedent
