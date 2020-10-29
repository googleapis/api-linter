---
rule:
  aip: 165
  name: [core, '0165', http-method]
  summary: Purge methods must use the POST HTTP verb.
permalink: /165/http-method
redirect_from:
  - /0165/http-method
---

# Purge methods: POST HTTP verb

This rule enforces that all `Purge` RPCs use the `POST` HTTP verb, as
mandated in [AIP-165][].

## Details

This rule looks at any message beginning with `Purge`, and complains
if the HTTP verb is anything other than `POST`. It _does_ check additional
bindings if they are present.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc PurgeBooks(PurgeBooksRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books:purge"  // Should be `post:`.
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
// (-- api-linter: core::0165::http-method=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc PurgeBooks(PurgeBooksRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books:purge"
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
