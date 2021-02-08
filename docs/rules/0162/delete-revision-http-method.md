---
rule:
  aip: 162
  name: [core, '0162', delete-revision-http-method]
  summary: Delete Revision methods must use the DELETE HTTP verb.
permalink: /162/delete-revision-http-method
redirect_from:
  - /0162/delete-revision-http-method
---

# Delete Revision methods: POST HTTP verb

This rule enforces that all Delete Revision RPCs use the `DELETE` HTTP verb, as
mandated in [AIP-162][].

## Details

This rule looks at any method matching `Delete*Revision`, and complains
if the HTTP verb is anything other than `DELETE`. It _does_ check additional
bindings if they are present.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc DeleteBookRevision(DeleteBookRevisionRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:deleteRevision"  // Should be `delete:`.
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
// (-- api-linter: core::0162::delete-revision-http-method=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc DeleteBookRevision(DeleteBookRevisionRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:deleteRevision"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent
