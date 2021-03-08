---
rule:
  aip: 162
  name: [core, '0162', list-revisions-http-method]
  summary: List Revisions methods must use the GET HTTP verb.
permalink: /162/list-revisions-http-method
redirect_from:
  - /0162/list-revisions-http-method
---

# List Revisions methods: GET HTTP verb

This rule enforces that all List Revisions RPCs use the `GET` HTTP verb, as
mandated in [AIP-162][].

## Details

This rule looks at any method matching `List*Revisions`, and complains
if the HTTP verb is anything other than `GET`. It _does_ check additional
bindings if they are present.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc ListBookRevisions(ListBookRevisionsRequest) returns (ListBookRevisionsResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:listRevisions"  // Should be `get:`.
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
// (-- api-linter: core::0162::list-revisions-http-method=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc ListBookRevisions(ListBookRevisionsRequest) returns (ListBookRevisionsResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:listRevisions"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent
