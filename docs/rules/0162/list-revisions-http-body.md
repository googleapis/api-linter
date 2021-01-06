---
rule:
  aip: 162
  name: [core, '0162', list-revisions-http-body]
  summary: List Revisions methods must not have an HTTP body.
permalink: /162/list-revisions-http-body
redirect_from:
  - /0162/list-revisions-http-body
---

# List Revisions methods: HTTP body

This rule enforces that all List Revisions RPCs omit the HTTP `body`,
as mandated in [AIP-162][].

## Details

This rule looks at any method matching `List*Revisions`, and complains
if the HTTP `body` field is set.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc ListBookRevisions(ListBookRevisionsRequest) returns (ListBookRevisionsResponse) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}:listRevisions"
    body: "*"  // This should be absent.
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc ListBookRevisions(ListBookRevisionsRequest) returns (ListBookRevisionsResponse) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}:listRevisions"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0162::list-revisions-http-body=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc ListBookRevisions(ListBookRevisionsRequest) returns (ListBookRevisionsResponse) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}:listRevisions"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent
