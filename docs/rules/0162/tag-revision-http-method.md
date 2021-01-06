---
rule:
  aip: 162
  name: [core, '0162', tag-revision-http-method]
  summary: Tag Revision methods must use the POST HTTP verb.
permalink: /162/tag-revision-http-method
redirect_from:
  - /0162/tag-revision-http-method
---

# Tag Revision methods: POST HTTP verb

This rule enforces that all Tag Revision RPCs use the `POST` HTTP verb, as
mandated in [AIP-162][].

## Details

This rule looks at any method matching `Tag*Revision`, and complains
if the HTTP verb is anything other than `POST`. It _does_ check additional
bindings if they are present.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc TagBookRevision(TagBookRevisionRequest) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}:tagRevision"  // Should be `post:`.
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
// (-- api-linter: core::0162::tag-revision-http-method=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc TagBookRevision(TagBookRevisionRequest) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}:tagRevision"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent
