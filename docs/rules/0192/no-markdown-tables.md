---
rule:
  aip: 192
  name: [core, '0192', no-markdown-tables]
  summary: Public comments should not include Markdown tables.
permalink: /192/no-markdown-tables
redirect_from:
  - /0192/no-markdown-tables
---

# No Markdown tables

This rule enforces that public comments for proto descriptors do not have
Markdown tables, as mandated in [AIP-192][].

## Details

This rule looks at each descriptor in each proto file (exempting oneofs and the
file itself) and complains if _public_ comments include Markdown tables.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Fields on the book include:
//
// Name     | Type
// -------- | --------
// `name`   | `string`
// `author` | `string`
message Book {
  // The resource name of the book.
  string name = 1;
}
```

**Correct** code for this rule:

```proto
// Correct.
// Fields on the book include:
//
// - `name`: `string`
// - `author`: `string`
message Book {
  // The resource name of the book.
  string name = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the descriptor.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// Fields on the book include:
//
// Name     | Type
// -------- | --------
// `name`   | `string`
// `author` | `string`
// (-- api-linter: core::0192::no-markdown-tables=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  // The resource name of the book.
  string name = 1;
}
```

[aip-192]: https://aip.dev/192
[aip.dev/not-precedent]: https://aip.dev/not-precedent
