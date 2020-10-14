---
rule:
  aip: 135
  name: [core, '0135', response-message-name]
  summary: Delete methods must return Empty or the resource.
permalink: /135/response-message-name
redirect_from:
  - /0135/response-message-name
---

# Delete methods: Response message

This rule enforces that all `Delete` RPCs have a response message of
`google.protobuf.Empty` or the resource, as mandated in [AIP-135][].

## Details

This rule looks at any message matching beginning with `Delete`, and complains
if the name of the corresponding output message is not one of:

- `google.protobuf.Empty`
- The name of the RPC with the prefix `Delete` removed.

**Important:** For declarative-friendly resources, only the resource is
permitted as a return type.

It also permits a response of `google.longrunning.Operation`; in this case, it
checks the `response_type` in the `google.longrunning.operation_info`
annotation and ensures that _it_ corresponds to either `google.protobuf.Empty`
or the name of the RPC with the prefix `Delete` removed.

## Examples

### Standard

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `google.protobuf.Empty` or the resource.
rpc DeleteBook(DeleteBookRequest) returns (DeleteBookResponse) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc DeleteBook(DeleteBookRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}"
  };
}
```

```proto
// Correct.
rpc DeleteBook(DeleteBookRequest) returns (Book) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}"
  };
}
```

**Important:** For declarative-friendly resources, only the resource is
permitted as a return type (and therefore only the second example is valid).

### Long-running operation

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc DeleteBook(DeleteBookRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}"
  };
  option (google.longrunning.operation_info) = {
    // Should be "google.protobuf.Empty" or "Book".
    response_type: "DeleteBookResponse"
    metadata_type: "DeleteBookMetadata"
  }
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc DeleteBook(DeleteBookRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}"
  };
  option (google.longrunning.operation_info) = {
    response_type: "google.protobuf.Empty"
    metadata_type: "DeleteBookMetadata"
  }
}
```

```proto
// Correct.
rpc DeleteBook(DeleteBookRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/books/*}"
  };
  option (google.longrunning.operation_info) = {
    response_type: "Book"
    metadata_type: "DeleteBookMetadata"
  }
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0135::response-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc DeleteBook(DeleteBookRequest) returns (DeleteBookResponse) {
  option (google.api.http) = {
    delete: "/v1/{book.name=publishers/*/books/*}"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-135]: https://aip.dev/135
[aip.dev/not-precedent]: https://aip.dev/not-precedent
