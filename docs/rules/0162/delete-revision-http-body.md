---
rule:
  aip: 162
  name: [core, '0162', delete-revision-http-body]
  summary: Delete Revision methods should not have an HTTP body.
permalink: /162/delete-revision-http-body
redirect_from:
  - /0162/delete-revision-http-body
---

# Delete Revision methods: HTTP body

This rule enforces that all Delete Revision RPCs have no HTTP `body`, as mandated in
[AIP-162][].

## Details

This rule looks at any method matching `Delete*Revision`, and complains
if it has an HTTP `body`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc DeleteBookRevision(DeleteBookRevisionRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}:deleteRevision"
    body: "*" // This should be removed.
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc DeleteBookRevision(DeleteBookRevisionRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}:deleteRevision"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0162::delete-revision-http-body=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc DeleteBookRevision(DeleteBookRevisionRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}:deleteRevision"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent
