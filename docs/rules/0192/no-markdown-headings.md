---
rule:
  aip: 192
  name: [core, '0192', no-markdown-headings]
  summary: Public comments should not include Markdown headings.
permalink: /192/no-markdown-headings
redirect_from:
  - /0192/no-markdown-headings
---

# No Markdown headings

This rule enforces that public comments for proto descriptors do not have
Markdown headings (`#`, `##`, etc.), as mandated in [AIP-192][].

## Details

This rule looks at each descriptor in each proto file (exempting oneofs and the
file itself) and complains if _public_ comments include Markdown headings (that
become `<h1>`, `<h2>`, etc.).

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// # A representation of a book.
message Book {
  // ## The resource name of the book.
  string name = 1;
}
```

**Correct** code for this rule:

```proto
// Correct.
// A representation of a book.
message Book {
  // The resource name of the book.
  string name = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the descriptor.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// # A representation of a book.
// (-- api-linter: core::0192::no-markdown-headings=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  // # The resource name of the book.
  // (-- api-linter: core::0192::no-markdown-headings=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string name = 1;
}
```

[aip-192]: https://aip.dev/192
[aip.dev/not-precedent]: https://aip.dev/not-precedent
