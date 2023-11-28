---
rule:
  aip: 192
  name: [core, '0192', closed-backticks]
  summary: Inline code should be surrounded by backticks.
permalink: /192/closed-backticks
redirect_from:
  - /0192/closed-backticks
---

# Closed backticks

This rule enforces that any inline code in public comments for proto descriptors
is surrounded by backticks, as mandated in [AIP-192][].

## Details

This rule looks at each descriptor in each proto file (exempting oneofs and the
file itself) and complains if _public_ comments include text that only has
either an opening backtick or a closing backtick.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Fields on the book include:
//
// - name`: `string
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
// - name`: `string
// (-- api-linter: core::0192::closed-backticks=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  // The resource name of the book.
  string name = 1;
}
```

[aip-192]: https://aip.dev/192
[aip.dev/not-precedent]: https://aip.dev/not-precedent
