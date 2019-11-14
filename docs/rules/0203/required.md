---
rule:
  aip: 203
  name: [core, '0203', required]
  summary: Required fields should be annotated.
permalink: /203/required
redirect_from:
  - /0203/required
---

# Required fields

This rule enforces that fields that are documented as required also have a
machine-readable annotation, as mandated by [AIP-203][].

## Details

This rule looks at any field with "required" (or similar forms) in the comment,
and complains if it does not have a `google.api.field_behavior` annotation.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;

  // @Required
  // The title of the book.
  string title = 2;
}
```

```proto
// Incorrect.
message Book {
  string name = 1;

  // Required. The title of the book.
  string title = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1;

  // The title of the book.
  string title = 2 [(google.api.field_behavior) = REQUIRED];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message Book {
  string name = 1;

  // Required. The title of the book.
  // (-- api-linter: core::0203::required
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string title = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-203]: https://aip.dev/203
[aip.dev/not-precedent]: https://aip.dev/not-precedent
