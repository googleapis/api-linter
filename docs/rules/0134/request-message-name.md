---
rule:
  aip: 134
  name: [core, '0134', request-message-name]
  summary: Update methods must have standardized request message names.
permalink: /134/request-message-name
redirect_from:
  - /0134/request-message-name
---

# Update methods: Request message

This rule enforces that all `Update` RPCs have a request message name of
`Update*Request`, as mandated in [AIP-134][].

## Details

This rule looks at any message matching beginning with `Update`, and complains
if the name of the corresponding input message does not match the name of the
RPC with the suffix `Request` appended.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc UpdateBook(Book) returns (Book) {  // Should be `UpdateBookRequest`.
  option (google.api.http) = {
    patch: "/v1/{name=publishers/*/books/*}"
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc UpdateBook(UpdateBookRequest) returns (Book) {
  option (google.api.http) = {
    patch: "/v1/{book.name=publishers/*/books/*}"
    body: "book"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0134::request-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc UpdateBook(Book) returns (Book) {  // Should be `UpdateBookRequest`.
  option (google.api.http) = {
    patch: "/v1/{name=publishers/*/books/*}"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-134]: https://aip.dev/134
[aip.dev/not-precedent]: https://aip.dev/not-precedent
