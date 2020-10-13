---
rule:
  aip: 134
  name: [core, '0134', response-lro]
  summary: |
    Declarative-friendly Update methods should use long-running operations.
permalink: /134/response-lro
redirect_from:
  - /0134/response-lro
---

# Long-running Update

This rule enforces that declarative-friendly update methods use long-running
operations, as mandated in [AIP-134][].

## Details

This rule looks at any `Update` method connected to a resource with a
`google.api.resource` annotation that includes `style: DECLARATIVE_FRIENDLY`,
and complains if it does not use long-running operations.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Assuming that Book is styled declarative-friendly, UpdateBook should
// return a long-running operation.
rpc UpdateBook(UpdateBookRequest) returns (Book) {
  option (google.api.http) = {
    patch: "/v1/{book.name=publishers/*/books/*}"
    body: "book"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
// Assuming that Book is styled declarative-friendly...
rpc UpdateBook(UpdateBookRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    patch: "/v1/{book.name=publishers/*/books/*}"
    body: "book"
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
// (-- api-linter: core::0134::response-lro=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc UpdateBook(UpdateBookRequest) returns (Book) {
  option (google.api.http) = {
    patch: "/v1/{book.name=publishers/*/books/*}"
    body: "book"
  };
}
```

**Note:** Violations of declarative-friendly rules should be rare, as tools are
likely to expect strong consistency.

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-134]: https://aip.dev/134
[aip.dev/not-precedent]: https://aip.dev/not-precedent
