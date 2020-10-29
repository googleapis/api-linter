---
rule:
  aip: 132
  name: [core, '0132', http-method]
  summary: List methods must use the GET HTTP verb.
permalink: /132/http-method
redirect_from:
  - /0132/http-method
---

# List methods: GET HTTP verb

This rule enforces that all `List` RPCs use the `GET` HTTP verb, as mandated in
[AIP-132][].

## Details

This rule looks at any message beginning with `List`, and complains if
the HTTP verb is anything other than `GET`. It _does_ check additional bindings
if they are present.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books"  // Should be `get:`.
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
// (-- api-linter: core::0132::http-method=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-132]: https://aip.dev/132
[aip.dev/not-precedent]: https://aip.dev/not-precedent
