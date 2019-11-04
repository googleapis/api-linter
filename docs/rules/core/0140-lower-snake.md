---
rule:
  aip: 140
  name: [core, '0140', lower-snake]
  summary: Field names should use `snake_case`.
---

# Field names: Abbreviations

This rule enforces that field names use `snake_case`, as mandated in
[AIP-140][].

## Details

This rule checks every field in the proto and complains if the field name
contains a capital letter.

## Examples

### Single word method

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;
  int32 pageCount = 2;  // Should be `page_count`.
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1;
  int32 page_count = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0140::lower-snake=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  string name = 1;
  string pageCount = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-140]: https://aip.dev/140
[aip.dev/not-precedent]: https://aip.dev/not-precedent
