---
rule:
  aip: 135
  name: [core, '0135', request-message-name]
  summary: Delete methods must have standardized request message names.
permalink: /135/request-message-name
redirect_from:
  - /0135/request-message-name
---

# Delete methods: Request message

This rule enforces that all `Delete` RPCs have a request message name of
`Delete*Request`, as mandated in [AIP-135][].

## Details

This rule looks at any message matching beginning with `Delete`, and complains
if the name of the corresponding input message does not match the name of the
RPC with the suffix `Request` appended.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `DeleteBookRequest`.
rpc DeleteBook(Book) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc DeleteBook(DeleteBookRequest) returns (Book) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0135::request-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc DeleteBook(Book) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-135]: https://aip.dev/135
[aip.dev/not-precedent]: https://aip.dev/not-precedent
