---
rule:
  aip: 165
  name: [core, '0165', request-message-name]
  summary: Purge methods must have standardized request message names.
permalink: /165/request-message-name
redirect_from:
  - /0165/request-message-name
---

# Purge methods: Request message

This rule enforces that all `Purge` RPCs have a request message name of
`Purge*Request`, as mandated in [AIP-165][].

## Details

This rule looks at any message beginning with `Purge`, and complains
if the name of the corresponding input message does not match the name of the
RPC with the suffix `Request` appended.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `PurgeBooksRequest`.
rpc PurgeBooks(DeleteBooksRequest) returns (google.longrunning.Operation) {
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
// (-- api-linter: core::0165::request-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc PurgeBooks(DeleteBooksRequest) returns (google.longrunning.Operation) {
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

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-165]: https://aip.dev/165
[aip.dev/not-precedent]: https://aip.dev/not-precedent
