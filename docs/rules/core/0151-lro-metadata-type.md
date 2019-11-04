---
rule:
  aip: 151
  name: [core, '0151', lro-metadata-type]
  summary: LRO methods must have a metadata type.
---

# LRO metadata type

This rule enforces that methods returning long-running operations specify a
metadata type in the `google.longrunning.operation_info` annotation, as
mandated in [AIP-151][].

## Details

This rule looks at any method with a return type of
`google.longrunning.Operation`, and complains if the `metadata_type` field
`google.longrunning.operation_info` annotation is not present.

Additionally, it complains if the metadata type is set to
`google.protobuf.Empty`, and recommends making an blank message instead, to
permit future extensibility.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc WriteBook(WriteBookRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books}:write"
    body: "*"
  };
  option (google.longrunning.operation_info) = {
    response_type: "WriteBookResponse"
    // `metadata_type` is not provided.
  };
}
```

```proto
// Incorrect.
rpc WriteBook(WriteBookRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books}:write"
    body: "*"
  };
  option (google.longrunning.operation_info) = {
    response_type: "WriteBookResponse"
    metadata_type: "google.protobuf.Empty"  // Should not be `Empty`.
  };
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
// (-- api-linter: core::0151::lro-metadata-type=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc WriteBook(WriteBookRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books}:write"
    body: "*"
  };
  option (google.longrunning.operation_info) = {
    response_type: "WriteBookResponse"
    metadata_type: "google.protobuf.Empty"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-151]: https://aip.dev/151
[aip.dev/not-precedent]: https://aip.dev/not-precedent
