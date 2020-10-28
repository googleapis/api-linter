---
rule:
  aip: 164
  name: [core, '0164', response-message-name]
  summary: Undelete methods must return the resource.
permalink: /164/response-message-name
redirect_from:
  - /0164/response-message-name
---

# Undelete methods: Response message

This rule enforces that all `Undelete` RPCs have a response message of
the resource, as mandated in [AIP-164][].

## Details

This rule looks at any message beginning with `Undelete`, and complains
if the name of the corresponding output message does not match the name of the
RPC with the prefix `Undelete` removed.

It also permits a response of `google.longrunning.Operation`; in this case, it
checks the `response_type` in the `google.longrunning.operation_info`
annotation and ensures that _it_ corresponds to the name of the RPC with the
prefix `Undelete` removed.

## Examples

### Standard

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `Book`.
rpc UndeleteBook(UndeleteBookRequest) returns (UndeleteBookResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:undelete"
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc UndeleteBook(UndeleteBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:undelete"
    body: "*"
  };
}
```

### Long-running operation

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc UndeleteBook(UndeleteBookRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:undelete"
    body: "*"
  };
  option (google.longrunning.operation_info) = {
    // Should be "Book".
    response_type: "UndeleteBookResponse"
    metadata_type: "UndeleteBookMetadata"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc UndeleteBook(UndeleteBookRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:undelete"
    body: "*"
  };
  option (google.longrunning.operation_info) = {
    response_type: "Book"
    metadata_type: "UndeleteBookMetadata"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0164::response-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc UndeleteBook(UndeleteBookRequest) returns (UndeleteBookResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:undelete"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-164]: https://aip.dev/164
[aip.dev/not-precedent]: https://aip.dev/not-precedent
