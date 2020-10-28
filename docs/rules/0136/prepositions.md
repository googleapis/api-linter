---
rule:
  aip: 136
  name: [core, '0136', prepositions]
  summary: Custom methods must not include prepositions in their names.
permalink: /136/prepositions
redirect_from:
  - /0136/prepositions
---

# Custom methods: Prepositions

This rule enforces that custom method names do not include most prepositions,
as mandated in [AIP-136][].

## Details

This rule looks at any method that is not a standard method, and complains if
it sees any of the following words in the method's name:

{% include prepositions.md %}

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// This RPC includes "with", which indicates a potential design concern.
rpc GetBookWithAuthor(GetBookWithAuthorRequest) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}:getWithAuthor"
  };
}
```

Sometimes, as in the example above, the inclusion of a preposition indicates a
design concern, and the method should be designed differently. In the above
case, the right answer is to stick to the standard method.

In other cases, the method may just need to be renamed, but with no major
conceptual change:

```proto
// Incorrect.
// The "FromLibrary" suffix is vestigial and should be removed.
rpc CheckoutBookFromLibrary(CheckoutBookFromLibraryRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:checkoutFromLibrary"
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc CheckoutBook(CheckoutBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:checkout"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0136::prepositions=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc GetBookWithAuthor(GetBookWithAuthorRequest) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}:getWithAuthor"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-136]: https://aip.dev/136
[aip.dev/not-precedent]: https://aip.dev/not-precedent
