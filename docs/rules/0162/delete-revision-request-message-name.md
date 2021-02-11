---
rule:
  aip: 162
  name: [core, '0162', delete-revision-request-message-name]
  summary: Delete Revision methods must have standardized request message names.
permalink: /162/delete-revision-request-message-name
redirect_from:
  - /0162/delete-revision-request-message-name
---

# Delete Revision methods: Request message

This rule enforces that all Delete Revision RPCs have a request message name of
`Delete*RevisionRequest`, as mandated in [AIP-162][].

## Details

This rule looks at any method matching `Delete*Revision`, and complains
if the name of the corresponding input message does not match the name of the
RPC with the suffix `Request` appended.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `DeleteBookRevisionRequest`.
rpc DeleteBookRevision(DeleteBookRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}:deleteRevision"
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
// (-- api-linter: core::0162::delete-revision-request-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc DeleteBookRevision(DeleteBookRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}:deleteRevision"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent
