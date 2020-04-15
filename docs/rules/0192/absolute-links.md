---
rule:
  aip: 192
  name: [core, '0192', absolute-links]
  summary: Comments must not use raw HTML
permalink: /192/absolute-links
redirect_from:
  - /0192/absolute-links
---

# Absolute links

This rule attempts to enforce that every descriptor in every proto file uses
absolute links, as mandated in [AIP-192][].

## Details

This rule looks at each descriptor in each proto file (exempting oneofs and the
file itself) and tries to find Markdown links using the `[link](uri)` syntax,
and complains if the URI does not have `://` in it.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// A representation of [a book](/wiki/Book).
message Book {
  string name = 1;
}
```

**Correct** code for this rule:

```proto
// Correct.
// A representation of [a book](https://en.wikipedia.org/wiki/Book).
message Book {
  string name = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the descriptor
(and revel in the irony). Remember to also include an [aip.dev/not-precedent][]
comment explaining why.

```proto
// (-- api-linter: core::0192::absolute-links=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
// A representation of [a book](/wiki/Book).
message Book {
  string name = 1;
}
```

[aip-192]: https://aip.dev/192
[aip.dev/not-precedent]: https://aip.dev/not-precedent
