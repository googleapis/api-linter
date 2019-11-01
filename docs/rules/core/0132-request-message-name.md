---
rule:
  aip: 132
  name: [core, '0132', request-message-name]
  summary: List methods must have standardized request message names.
---

# List methods: Request message

This rule enforces that all `List` RPCs have a request message name of
`List*Request`, as mandated in [AIP-132][].

## Details

This rule looks at any message matching beginning with `List`, and complains if
the name of the corresponding input message does not match the name of the RPC
with the suffix `Request` appended.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `ListBooksRequest`.
rpc ListBooks(ListBooksReq) returns (ListBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books"
  }
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books"
  }
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0132::request-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc ListBooks(ListBooksReq) returns (ListBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books"
  }
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-132]: https://aip.dev/132
[aip.dev/not-precedent]: https://aip.dev/not-precedent
