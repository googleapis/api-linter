---
rule:
  aip: 192
  name: [core, '0192', no-html]
  summary: Comments must not use raw HTML
permalink: /192/no-html
redirect_from:
  - /0192/no-html
---

# No HTML in comments

This rule enforces that every descriptor in every proto file does not use raw
HTML in comments, as mandated in [AIP-192][].

## Details

This rule looks at each descriptor in each proto file (exempting oneofs and the
file itself) and tries to pick up "HTML smell", and complains if it thinks it
found HTML.

**Note:** This lint rule uses a regular expression to look for HTML, which is a
[famous anti-pattern][]. We do it anyway to avoid taking a large dependency for
one rule. Therefore, this rule allows many false negatives to avoid any false
positives; that is, it will intentionally let more complex HTML through in
order to prevent cases where it complains and is wrong.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// A representation of a book.
message Book {
  // (-- This comment should use Markdown, not HTML.) --)
  // The name of the book.
  // Format: <code>publishers/{publisher}/books/{book}</code>
  string name = 1;
}
```

**Correct** code for this rule:

```proto
// Correct.
// A representation of a book.
// A representation of a book.
message Book {
  // The name of the book.
  // Format: `publishers/{publisher}/books/{book}`
  string name = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the descriptor
(and revel in the irony). Remember to also include an [aip.dev/not-precedent][]
comment explaining why.

```proto
// A representation of a book.
message Book {
  // (-- api-linter: core::0192::no-html=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  // The name of the book.
  // Format: <code>publishers/{publisher}/books/{book}</code>
  string name = 1;
}
```

[aip-192]: https://aip.dev/192
[aip.dev/not-precedent]: https://aip.dev/not-precedent
[famous anti-pattern]: https://stackoverflow.com/questions/1732348/
