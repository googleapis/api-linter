---
rule:
  aip: 162
  name: [core, '0162', tag-revision-request-message-name]
  summary: Tag Revision methods must have standardized request message names.
permalink: /162/tag-revision-request-message-name
redirect_from:
  - /0162/tag-revision-request-message-name
---

# Tag Revision methods: Request message

This rule enforces that all Tag Revision RPCs have a request message name of
`Tag*RevisionRequest`, as mandated in [AIP-162][].

## Details

This rule looks at any method matching `Tag*Revision`, and complains
if the name of the corresponding input message does not match the name of the
RPC with the suffix `Request` appended.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `TagBookRevisionRequest`.
rpc TagBookRevision(TagBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:tagRevision"
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc TagBookRevision(TagBookRevisionRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:tagRevision"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0162::tag-revision-request-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc TagBookRevision(TagBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:tagRevision"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent
