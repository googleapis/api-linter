---
rule:
  aip: 134
  name: [core, '0134', synonyms]
  summary: Update methods must be named starting with "Update".
permalink: /134/synonyms
redirect_from:
  - /0134/synonyms
---

# Update methods: Synonym check

This rule enforces that single-resource creation methods have names beginning
with `Update`, as mandated in [AIP-134][].

## Details

This rule looks at any message with names similar to `Update`, and suggests
using `Update` instead. It complains if it sees the following synonyms:

- Patch
- Put
- Set

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc PatchBook(PatchBookRequest) returns (Book) {  // Should be `UpdateBook`.
  option (google.api.http) = {
    patch: "/v1/{book.name=publishers/*/books/*}"
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
// (-- api-linter: core::0134::synonyms=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc PatchBook(PatchBookRequest) returns (Book) {
  option (google.api.http) = {
    patch: "/v1/{book.name=publishers/*/books/*}"
    body: "book"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-134]: https://aip.dev/134
[aip.dev/not-precedent]: https://aip.dev/not-precedent
