---
rule:
  aip: 133
  name: [core, '0133', response-lro]
  summary: |
    Declarative-friendly create methods should use long-running operations.
permalink: /133/response-lro
redirect_from:
  - /0133/response-lro
---

# Long-running create

This rule enforces that declarative-friendly create methods use long-running
operations, as mandated in [AIP-133][].

## Details

This rule looks at any `Create` method connected to a resource with a
`google.api.resource` annotation that includes `style: DECLARATIVE_FRIENDLY`,
and complains if it does not use long-running operations.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Assuming that Book is styled declarative-friendly, CreateBook should
// return a long-running operation.
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books"
    body: "book"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
// Assuming that Book is styled declarative-friendly...
rpc CreateBook(CreateBookRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books"
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
// (-- api-linter: core::0133::response-lro=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books"
    body: "book"
  };
}
```

**Note:** Violations of declarative-friendly rules should be rare, as tools are
likely to expect strong consistency.

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-133]: https://aip.dev/133
[aip.dev/not-precedent]: https://aip.dev/not-precedent
