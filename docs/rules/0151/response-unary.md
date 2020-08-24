---
rule:
  aip: 151
  name: [core, '0151', response-unary]
  summary: Long-running operations must not use streaming.
permalink: /151/response-unary
redirect_from:
  - /0151/response-unary
---

# Paginated methods: Unary responses

This rule enforces that all long-running operation methods use unary responses,
as mandated in [AIP-151][].

## Details

This rule looks at any message returning a `google.longrunning.Operation`, and
complains if the method uses gRPC server streaming (the `stream` keyword).

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Streaming is prohibited on long-running operations.
rpc ReadBook(ReadBookRequest) returns (stream google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*/books/*}:read"
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc ReadBook(ReadBookRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*/books/*}:read"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message or
above the field. Remember to also include an [aip.dev/not-precedent][] comment
explaining why.

```proto
// (-- api-linter: core::0151::response-unary
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc ReadBook(ReadBookRequest) returns (stream google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*/books/*}:read"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-151]: https://aip.dev/151
[aip.dev/not-precedent]: https://aip.dev/not-precedent
