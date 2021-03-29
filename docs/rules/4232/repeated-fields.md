---
rule:
  aip: 4232
  name: [client-libraries, '4232', repeated-fields]
  summary: Method Signatures can only have repeated fields in the last element.
permalink: /4232/repeated-fields
redirect_from:
  - /4232/repeated-fields
---

# Method Signature: Repeated fields

This rule enforces that all `google.api.method_signature` annotations do not
have `repeated` fields in any position other than the last field of a signature,
as mandated in [AIP-4232][].

## Details

This rule looks at any RPC methods with a `google.api.method_signature`
annotation, and complains if any field other than the last is `repeated`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc BatchDeleteBooks(BatchDeleteBooksRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchDelete"
    body: "*"
  };
  // The field "names" is repeated and must only come at the end.
  option (google.api.method_signature) = "names,parent"
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc BatchDeleteBooks(BatchDeleteBooksRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchDelete"
    body: "*"
  };
  option (google.api.method_signature) = "parent,names"
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: client-libraries::4232::repeated-fields=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc BatchDeleteBooks(BatchDeleteBooksRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books:batchDelete"
    body: "*"
  };
  option (google.api.method_signature) = "names,parent"
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-4232]: https://aip.dev/4232
[aip.dev/not-precedent]: https://aip.dev/not-precedent
