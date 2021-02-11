---
rule:
  aip: 162
  name: [core, '0162', delete-revision-http-uri-suffix]
  summary: Delete Revision methods must have the correct URI suffix
permalink: /162/delete-revision-http-uri-suffix
redirect_from:
  - /0162/delete-revision-http-uri-suffix
---

# Delete Revision methods: URI suffix

This rule enforces that Delete Revision methods include the `:deleteRevision` suffix
in the REST URI, as mandated in [AIP-162][].

## Details

This rule looks at any method matching `Delete*Revision`, and
complains if the HTTP URI does not end with `:deleteRevision`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc DeleteBookRevision(DeleteBookRevisionRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}:delete"  // Should end with `:deleteRevision`
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
// (-- api-linter: core::0162::delete-revision-http-uri-suffix=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc DeleteBookRevision(DeleteBookRevisionRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}:delete"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent
