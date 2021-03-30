---
rule:
  aip: 132
  name: [core, '0132', http-uri-parent]
  summary: List methods must map the parent field to the URI.
permalink: /132/http-uri-parent
redirect_from:
  - /0132/http-uri-parent
---

# List methods: HTTP URI parent field

This rule enforces that all `List` RPCs map the `parent` field to the HTTP
URI, as mandated in [AIP-132][].

## Details

This rule looks at any message beginning with `List`, and complains
if the `parent` variable is not included in the URI. It _does_ check additional
bindings if they are present.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {
  option (google.api.http) = {
    get: "/v1/publishers/*/books"  // The `parent` field should be extracted.
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
// (-- api-linter: core::0132::http-uri-parent=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {
  option (google.api.http) = {
    get: "/v1/publishers/*/books"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-132]: https://aip.dev/132
[aip.dev/not-precedent]: https://aip.dev/not-precedent
