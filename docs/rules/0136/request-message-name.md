---
rule:
  aip: 136
  name: [core, '0136', request-message-name]
  summary: Custom methods must have standardized request message names.
permalink: /136/request-message-name
redirect_from:
  - /0136/request-message-name
---

# Custom methods: Request message

This rule enforces that all custom methods should take a request message
matching the RPC name, with a `Request` suffix [AIP-136][].

## Details

This rule looks at any method that is not a standard method, and complains if
the name of the corresponding input message does not match the name of the RPC
with the suffix `Request` appended.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `ArchiveBookRequest`.
rpc ArchiveBook(Book) returns (ArchiveBookResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:archive"
    body: "*"
  };
}

```

**Correct** code for this rule:

```proto
// Correct.
rpc ArchiveBook(ArchiveBookRequest) returns (ArchiveBookResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:archive"
    body: "*"
  };
}

```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0136::request-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc ArchiveBook(Book) returns (ArchiveBookResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:archive"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-136]: https://aip.dev/136
[aip.dev/not-precedent]: https://aip.dev/not-precedent
