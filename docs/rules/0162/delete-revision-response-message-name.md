---
rule:
  aip: 162
  name: [core, '0162', delete-revision-response-message-name]
  summary: Delete Revision methods must return the resource.
permalink: /162/delete-revision-response-message-name
redirect_from:
  - /0162/delete-revision-response-message-name
---

# Delete Revision methods: Response message

This rule enforces that all Delete Revision RPCs have a response message of
the resource, as mandated in [AIP-162][].

## Details

This rule looks at any message matching `Delete*Revision`, and complains
if the corresponding output message does not match the name of the RPC with the
prefix `Delete` and suffix `Revision` removed.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should return `Book`.
rpc DeleteBookRevision(DeleteBookRevisionRequest) returns (DeleteBookRevisionResponse) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}:deleteRevision"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc DeleteBookRevision(DeleteBookRevisionRequest) returns (Book) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}:deleteRevision"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0162::delete-revision-response-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc DeleteBookRevision(DeleteBookRevisionRequest) returns (DeleteBookRevisionResponse) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}:deleteRevision"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent
