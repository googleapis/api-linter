---
rule:
  aip: 216
  name: [core, '0216', synonyms]
  summary: Lifecycle state enums should be called "State", not "Status".
permalink: /216/synonyms
redirect_from:
  - /0216/synonyms
---

# States

This rule enforces that all lifecycle state enums are called `State` rather
than `Status`, as mandated in [AIP-216][].

## Details

This rule iterates over enumerations and looks for enums with a name of
`Status` or ending in `Status`, and suggests the use of `State` instead.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
enum Status {  // Should be `State`.
  STATUS_UNSPECIFIED = 0;
}
```

```proto
// Incorrect.
enum BookStatus {  // Should be `Book.State` or `BookState`.
  BOOK_STATUS_UNSPECIFIED = 0;
  HARDCOVER = 1;
}
```

**Correct** code for this rule:

```proto
// Correct.
enum State {
  STATE_UNSPECIFIED = 0;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the enum value.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0216::synonyms=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
enum Status {
  STATUS_UNSPECIFIED = 0;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-216]: https://aip.dev/216
[aip.dev/not-precedent]: https://aip.dev/not-precedent
