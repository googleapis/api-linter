---
rule:
  aip: 133
  name: [core, '0133', synonyms]
  summary: Create methods must be named starting with "Create".
permalink: /133/synonyms
redirect_from:
  - /0133/synonyms
---

# Create methods: Synonym check

This rule enforces that single-resource creation methods have names beginning
with `Create`, as mandated in [AIP-133][].

## Details

This rule looks at any message with names similar to `Create`, and suggests
using `Create` instead. It complains if it sees the following synonyms:

- Insert
- Make
- Post

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc InsertBook(InsertBookRequest) returns (Book) {  // Should be `CreateBook`.
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books"
    body: "book"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books"
    body: "book"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0133::synonyms=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc InsertBook(InsertBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books"
    body: "book"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-133]: https://aip.dev/133
[aip.dev/not-precedent]: https://aip.dev/not-precedent
