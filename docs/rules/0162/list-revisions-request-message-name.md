---
rule:
  aip: 162
  name: [core, '0162', list-revisions-request-message-name]
  summary: List Revisions methods must have standardized request message names.
permalink: /162/list-revisions-request-message-name
redirect_from:
  - /0162/list-revisions-request-message-name
---

# List Revisions methods: Request message

This rule enforces that all List Revisions RPCs have a request message name of
`List*RevisionsRequest`, as mandated in [AIP-162][].

## Details

This rule looks at any method matching `List*Revisions`, and complains
if the name of the corresponding input message does not match the name of the
RPC with the suffix `Request` appended.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `ListBookRevisionsRequest`.
rpc ListBookRevisions(ListRevisionsRequest) returns (ListBookRevisionsResponse) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}:listRevisions"
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
// (-- api-linter: core::0162::list-revisions-request-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc ListBookRevisions(ListRevisionsRequest) returns (ListBookRevisionsResponse) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}:listRevisions"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent
