---
rule:
  aip: 162
  name: [core, '0162', list-revisions-http-uri-suffix]
  summary: List Revisions methods must have the correct URI suffix
permalink: /162/list-revisions-http-uri-suffix
redirect_from:
  - /0162/list-revisions-http-uri-suffix
---

# List Revisions methods: URI suffix

This rule enforces that List Revisions methods include the `:listRevisions` suffix
in the REST URI, as mandated in [AIP-162][].

## Details

This rule looks at any method matching `List*Revisions`, and
complains if the HTTP URI does not end with `:listRevisions`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc ListBookRevisions(ListBookRevisionsRequest) returns (ListBookRevisionsResponse) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}:list"  // Should end with `:listRevisions`
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
// (-- api-linter: core::0162::list-revisions-http-uri-suffix=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc ListBookRevisions(ListBookRevisionsRequest) returns (ListBookRevisionsResponse) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}:list"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent
