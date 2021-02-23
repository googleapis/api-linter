---
rule:
  aip: 162
  name: [core, '0162', tag-revision-http-uri-suffix]
  summary: Tag Revision methods must have the correct URI suffix
permalink: /162/tag-revision-http-uri-suffix
redirect_from:
  - /0162/tag-revision-http-uri-suffix
---

# Tag Revision methods: URI suffix

This rule enforces that Tag Revision methods include the `:tagRevision` suffix
in the REST URI, as mandated in [AIP-162][].

## Details

This rule looks at any method matching `Tag*Revision`, and
complains if the HTTP URI does not end with `:tagRevision`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc TagBookRevision(TagBookRevisionRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:tag"  // Should end with `:tagRevision`
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
// (-- api-linter: core::0162::tag-revision-http-uri-suffix=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc TagBookRevision(TagBookRevisionRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:tag"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent
