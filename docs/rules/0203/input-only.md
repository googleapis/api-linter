---
rule:
  aip: 203
  name: [core, '0203', input-only]
  summary: Input only fields should be annotated.
permalink: /203/input-only
redirect_from:
  - /0203/input-only
---

# Input only fields

This rule enforces that fields that are documented as input only also have a
machine-readable annotation, as mandated by [AIP-203][].

## Details

This rule looks at any field with "input only" (or similar forms) in the
comment, and complains if it does not have a `google.api.field_behavior`
annotation.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;

  // @InputOnly
  // The password used to check out this book.
  string access_password = 2;
}
```

```proto
// Incorrect.
message Book {
  string name = 1;

  // Input only. The password used to check out this book.
  string access_password = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1;

  // The password used to check out this book.
  string access_password = 2 [(google.api.field_behavior) = INPUT_ONLY];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message Book {
  string name = 1;

  // Input only. The password used to check out this book.
  // (-- api-linter: core::0203::input-only=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string access_password = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-203]: https://aip.dev/203
[aip.dev/not-precedent]: https://aip.dev/not-precedent
