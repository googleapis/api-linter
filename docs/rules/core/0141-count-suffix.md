---
rule:
  aip: 141
  name: [core, '0141', count-suffix]
  summary: Quantities should use a `_count` suffix.
---

# Count suffix

This rule tries to enforce that discrete quantities have consistent field names
ending in `_count`, as mandated in [AIP-141][].

## Details

This rule scans all fields and complains if it sees a `num_` prefix, and
suggests a `_count` suffix instead.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;
  int32 num_pages = 2;  // Should be `page_count`.
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
// (-- api-linter: core::0141::count-suffix=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  string name = 1;
  int32 num_pages = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-141]: https://aip.dev/141
[aip.dev/not-precedent]: https://aip.dev/not-precedent
