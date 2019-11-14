---
rule:
  aip: 132
  name: [core, '0132', response-message-name]
  summary: List methods must have standardized response message names.
permalink: /132/response-message-name
redirect_from:
  - /0132/response-message-name
---

# List methods: Response message

This rule enforces that all `List` RPCs have a response message name of
`List*Response`, as mandated in [AIP-132][].

## Details

This rule looks at any message matching beginning with `List`, and complains if
the name of the corresponding returned message does not match the name of the
RPC with the suffix `Response` appended.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `ListBooksResponse`.
rpc ListBooks(ListBooksRequest) returns (Books) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0132::response-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc ListBooks(ListBooksRequest) returns (Books) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-132]: https://aip.dev/132
[aip.dev/not-precedent]: https://aip.dev/not-precedent
