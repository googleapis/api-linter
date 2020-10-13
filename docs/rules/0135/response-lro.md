---
rule:
  aip: 135
  name: [core, '0135', response-lro]
  summary: |
    Declarative-friendly delete methods should use long-running operations.
permalink: /135/response-lro
redirect_from:
  - /0135/response-lro
---

# Long-running Delete

This rule enforces that declarative-friendly delete methods use long-running
operations, as mandated in [AIP-135][].

## Details

This rule looks at any `Delete` method connected to a resource with a
`google.api.resource` annotation that includes `style: DECLARATIVE_FRIENDLY`,
and complains if it does not use long-running operations.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Assuming that Book is styled declarative-friendly, DeleteBook should
// return a long-running operation.
rpc DeleteBook(DeleteBookRequest) returns (Book) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
// Assuming that Book is styled declarative-friendly...
rpc DeleteBook(DeleteBookRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}"
  };
  option (google.longrunning.operation_info) = {
    response_type: "Book"
    metadata_type: "OperationMetadata"
  }
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0135::response-lro=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc DeleteBook(DeleteBookRequest) returns (Book) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}"
    body: "book"
  };
}
```

**Note:** Violations of declarative-friendly rules should be rare, as tools are
likely to expect strong consistency.

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-135]: https://aip.dev/135
[aip.dev/not-precedent]: https://aip.dev/not-precedent
