---
rule:
  aip: 192
  name: [core, '0192', deprecated-comment]
  summary: Deprecated elements must have a corresponding comment.
permalink: /192/deprecated-comment
redirect_from:
  - /0192/deprecated-comment
---

# Deprecated comments

This rule enforces that every element marked with the protobuf `deprecated`
option has `"Deprecated: <reason>"` as the first line in the public leading
comment, as mandated in [AIP-192][].

## Details

This rule looks at each descriptor in each proto file, and complains if the
protobuf `deprecated` option is set to `true`, but the first line of the public
comment does not begin with "Deprecated: ".

## Examples

**Incorrect** code for this rule:

```proto
// A library service.
service Library {
  // Incorrect.
  // Retrieves a book.
  rpc GetBook(GetBookRequest) returns (Book) {
    option deprecated = true;
  }
}
```

**Correct** code for this rule:

```proto
// A library service.
service Library {
  // Deprecated: Please borrow a physical book instead.
  // Correct.
  // Retrieves a book.
  rpc GetBook(GetBookRequest) returns (Book) {
    option deprecated = true;
  }
}
```

## Disabling

If you need to violate this rule, use a leading comment above the descriptor.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// A library service.
service Library {
  // (-- api-linter: core::0192::deprecated-comment=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  // Incorrect.
  // Retrieves a book.
  rpc GetBook(GetBookRequest) returns (Book) {
    option deprecated = true;
  }
}
```

[aip-192]: https://aip.dev/192
[aip.dev/not-precedent]: https://aip.dev/not-precedent
