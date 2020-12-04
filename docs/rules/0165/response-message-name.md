---
rule:
  aip: 165
  name: [core, '0165', response-message-name]
  summary: Purge methods must return a long-running operation.
permalink: /165/response-message-name
redirect_from:
  - /0165/response-message-name
---

# Purge methods: Response message

This rule enforces that all `Purge` RPCs return a long-running operation that
resolves to a response with a corresponding name, as mandated in [AIP-165][].

## Details

This rule looks at any message beginning with `Purge`, and complains if the
response is not a long-running operation that resolves to a response matching
the name of the method with a `Response` suffix.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `google.longrunning.Operation`.
rpc PurgeBooks(PurgeBooksRequest) returns (PurgeBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:purge"
    body: "*"
  };
}
```

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc PurgeBooks(PurgeBooksRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:purge"
    body: "*"
  };
  option (google.longrunning.operation_info) = {
    // Should be "PurgeBooksResponse".
    response_type: "Book"
    metadata_type: "PurgeBooksMetadata"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc PurgeBooks(PurgeBooksRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:purge"
    body: "*"
  };
  option (google.longrunning.operation_info) = {
    response_type: "PurgeBooksResponse"
    metadata_type: "PurgeBooksMetadata"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0165::response-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc PurgeBooks(PurgeBooksRequest) returns (PurgeBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:purge"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-165]: https://aip.dev/165
[aip.dev/not-precedent]: https://aip.dev/not-precedent
