---
rule:
  aip: 133
  name: [core, '0133', request-message-name]
  summary: Create methods must have standardized request message names.
permalink: /133/request-message-name
redirect_from:
  - /0133/request-message-name
---

# Create methods: Request message

This rule enforces that all `Create` RPCs have a request message name of
`Create*Request`, as mandated in [AIP-133][].

## Details

This rule looks at any message matching beginning with `Create`, and complains
if the name of the corresponding input message does not match the name of the
RPC with the suffix `Request` appended.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc CreateBook(Book) returns (Book) {  // Should be `CreateBookRequest`.
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books"
    body: "*"
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
// (-- api-linter: core::0133::request-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc CreateBook(Book) returns (Book) {
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
