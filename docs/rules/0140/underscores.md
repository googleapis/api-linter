---
rule:
  aip: 140
  name: [core, '0140', underscores]
  summary: Field names must not have goofy underscores.
permalink: /140/underscores
redirect_from:
  - /0140/underscores
---

# Field names: Underscores

This rule enforces that field names do not use leading, trailing, or adjacent
underscores, as mandated in [AIP-140][].

## Details

This rule checks every field in the proto and complains if it sees leading or
trailing underscores, or two or more underscores with nothing in between them.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;
  string _title = 2;  // Should be `title`.
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1;
  string title = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

**Warning:** Violating this rule is likely to run into tooling failures.

```proto
// (-- api-linter: core::0140::underscores=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  string name = 1;
  string _title = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-140]: https://aip.dev/140
[aip.dev/not-precedent]: https://aip.dev/not-precedent
