---
rule:
  aip: 164
  name: [core, '0164', response-lro]
  summary: |
    Declarative-friendly undelete methods should use long-running operations.
permalink: /164/response-lro
redirect_from:
  - /0164/response-lro
---

# Long-running Undelete

This rule enforces that declarative-friendly undelete methods use long-running
operations, as mandated in [AIP-164][].

## Details

This rule looks at any `Undelete` method connected to a resource with a
`google.api.resource` annotation that includes `style: DECLARATIVE_FRIENDLY`,
and complains if it does not use long-running operations.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Assuming that Book is styled declarative-friendly, UndeleteBook should
// return a long-running operation.
rpc UndeleteBook(UndeleteBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:undelete"
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
// Assuming that Book is styled declarative-friendly...
rpc UndeleteBook(UndeleteBookRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:undelete"
    body: "*"
  };
  option (google.longrunning.operation_info) = {
    response_type: "Book"
    metadata_type: "OperationMetadata"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0164::response-lro=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc UndeleteBook(UndeleteBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:undelete"
    body: "*"
  };
}
```

**Note:** Violations of declarative-friendly rules should be rare, as tools are
likely to expect strong consistency.

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-164]: https://aip.dev/164
[aip.dev/not-precedent]: https://aip.dev/not-precedent
