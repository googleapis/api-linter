---
rule:
  aip: 126
  name: [core, '0126', upper-snake-values]
  summary: All enum values must be in upper snake case.
permalink: /126/upper-snake-values
redirect_from:
  - /0126/upper-snake-values
---

# Upper snake case values

This rule enforces that all enum values be in upper snake case, as mandated in
[AIP-126][].

## Details

This rule finds all enumerations and ensures that each value is provided in
`UPPER_SNAKE_CASE`: uppercase letters and numbers, with words separated by a
single underscore. Leading and trailing underscores are not allowed.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
enum Format {
  FORMAT_UNSPECIFIED = 0;
  hardcover = 1;  // Should be "HARDCOVER".
  _FOO = 2;       // Should be "FOO".
  BAR_ = 3;       // Should be "BAR".
}
```

**Correct** code for this rule:

```proto
// Correct.
enum Format {
  FORMAT_UNSPECIFIED = 0;
  HARDCOVER = 1;
  FOO = 2;
  BAR = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the enum value.

```proto
enum Format {
  FORMAT_UNSPECIFIED = 0;

  // (-- api-linter: core::0126::upper-snake-values=disabled --)
  hardcover = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-126]: https://aip.dev/126
