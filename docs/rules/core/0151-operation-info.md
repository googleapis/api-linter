---
rule:
  aip: 151
  name: [core, '0151', operation-info]
  summary: LRO methods must include an `operation_info` annotation.
---

# Long-running operation info

This rule enforces that methods returning long-running operations include an
annotation specifying their response type and metadata type, as mandated in
[AIP-151][].

## Details

This rule looks at any method with a return type of
`google.longrunning.Operation`, and complains if the
`google.longrunning.operation_info` annotation is not present.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc WriteBook(WriteBookRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books}:write"
    body: "*"
  };
  // There should be a google.longrunning.operation_info annotation.
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc WriteBook(WriteBookRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books}:write"
    body: "*"
  };
  option (google.longrunning.operation_info) = {
    response_type: "WriteBookResponse"
    metadata_type: "WriteBookMetadata"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0151::operation-info=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc WriteBook(WriteBookRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books}:write"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-151]: https://aip.dev/151
[aip.dev/not-precedent]: https://aip.dev/not-precedent
