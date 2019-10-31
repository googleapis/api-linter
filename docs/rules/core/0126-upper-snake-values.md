---
rule:
  aip: 126
  name: [core, '0126', upper-snake-values]
  summary: All enum values must be in upper snake case.
---

# Upper snake case values

This rule enforces that all enum values be in upper snake case, as mandated in
[AIP-126](http://aip.dev/126).

## Details

This rule finds all enumerations and ensures that each value is provided in
`UPPER_SNAKE_CASE`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
enum Format {
  FORMAT_UNSPECIFIED = 0;
  hardcover = 1;  // Should be "HARDCOVER".
}
```

**Correct** code for this rule:

```proto
// Correct.
enum Format {
  FORMAT_UNSPECIFIED = 0;
  HARDCOVER = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the enum value.

```proto
enum Format {
  FORMAT_UNSPECIFIED = 0;

  // (-- api-linter: core::0126::upper-snake-values=disabled --)
  hardcover = 0;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.
