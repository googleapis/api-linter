---
rule:
  aip: 134
  name: [core, '0134', http-method]
  summary: Update methods must use the PATCH HTTP verb.
permalink: /134/http-method
redirect_from:
  - /0134/http-method
---

# Update methods: PATCH HTTP verb

This rule enforces that all `Update` RPCs use the `PATCH` HTTP verb, as
mandated in [AIP-134][].

## Details

This rule looks at any message matching beginning with `Update`, and complains
if the HTTP verb is anything other than `PATCH`. It _does_ check additional
bindings if they are present.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc UpdateBook(UpdateBookRequest) returns (Book) {
  option (google.api.http) = {
    put: "/v1/{book.name=publishers/*/books/*}"  // Should be `patch:`.
    body: "book"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc UpdateBook(UpdateBookRequest) returns (Book) {
  option (google.api.http) = {
    patch: "/v1/{book.name=publishers/*/books/*}"
    body: "book"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0134::http-method=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc UpdateBook(UpdateBookRequest) returns (Book) {
  option (google.api.http) = {
    put: "/v1/{book.name=publishers/*/books/*}"  // Should be `patch:`.
    body: "book"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-134]: https://aip.dev/134
[aip.dev/not-precedent]: https://aip.dev/not-precedent
