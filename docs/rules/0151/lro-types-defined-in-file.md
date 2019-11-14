---
rule:
  aip: 151
  name: [core, '0151', lro-types-defined-in-file]
  summary: LRO response and metadata messages must be defined in the same file.
permalink: /151/lro-types-defined-in-file
redirect_from:
  - /0151/lro-types-defined-in-file
---

# LRO types defined in file

This rule enforces that methods returning long-running operations define their
response and metadata messages in the same file, as mandated in [AIP-151][].

## Details

This rule looks at any method with a return type of
`google.longrunning.Operation`, and complains if the messages designated by the
`response_type` and `metadata_type` fields are not defined in the same file.

Because these message names are strings, and a string reference does not
require an `import` statement, defining the response and metadata types
elsewhere can cause problems for tooling. To prevent this, and also to maintain
consistency with the file layout in [AIP-191][], the linter complains if the
message is not defined in the same file.

## Examples

**Incorrect** code for this rule:

In `library_service.proto`:

```proto
// Incorrect.
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

In `operations.proto`:

```proto
// Incorrect.
message WriteBookResponse {
  // Should be in the same file.
}

message WriteBookMetadata {
  // Should be in the same file.
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

// Later in the file...
message WriteBookResponse {
  // ...
}

message WriteBookMetadata {
  // ...
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0151::lro-types-defined-in-file=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc WriteBook(WriteBookRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books}:write"
    body: "*"
  };
  option (google.longrunning.operation_info) = {
    response_type: "google.protobuf.Empty"
    metadata_type: "WriteBookMetadata"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-151]: https://aip.dev/151
[aip-191]: https://aip.dev/191
[aip.dev/not-precedent]: https://aip.dev/not-precedent
