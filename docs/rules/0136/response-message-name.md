---
rule:
  aip: 136
  name: [core, '0136', response-message-name]
  summary: Custom methods must have standardized response message names.
permalink: /136/response-message-name
redirect_from:
  - /0136/response-message-name
---

# Custom methods: Response message

This rule enforces that all custom methods should take a response message
matching the RPC name, with a `Response` suffix, or the resource being operated
on [AIP-136][].

## Details

This rule looks at any method that is not a standard method, and complains if
the name of the corresponding input message does not match the name of the RPC
with the suffix `Response` appended, or the resource itself. It will use the
`(google.api.resource_reference)` on the first field of the request message to
determine the resource name that should be used.

## Examples

### Response Suffix

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `ArchiveBookResponse`.
rpc ArchiveBook(ArchiveBookRequest) returns (ArchiveBookResp) {
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

```

### Resource

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `Book`.
rpc ArchiveBook(ArchiveBookRequest) returns (Author) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:archive"
    body: "*"
  };
}

message ArchiveBookRequest {
    // The book to archive.
    // Format: publishers/{publisher}/books/{book}
    string name = 1 [(google.api.resource_reference).type = "library.googleapis.com/Book"];
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc ArchiveBook(ArchiveBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:archive"
    body: "*"
  };
}

message ArchiveBookRequest {
    // The book to archive.
    // Format: publishers/{publisher}/books/{book}
    string name = 1 [(google.api.resource_reference).type = "library.googleapis.com/Book"];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0136::response-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc ArchiveBook(ArchiveBookRequest) returns (ArchiveBookResp) {
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
