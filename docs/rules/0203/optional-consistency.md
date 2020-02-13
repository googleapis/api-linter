---
rule:
  aip: 203
  name: [core, '0203', optional-consistency]
  summary: Optional fields should either be always or never annotated.
permalink: /203/optional-consistency
redirect_from:
  - /0203/optional-consistency
---

# Optional fields: Consistency

This rule enforces that messages containing fields that have a
`google.api.field_behavior` annotation of `OPTIONAL` uses this consistently
throughout the message, as mandated by [AIP-203][].

## Details

This rule looks at messages with at least one field with a
`google.api.field_behavior` annotation of `OPTIONAL`, and complains if it finds
other fields on the message with no `google.api.field_behavior` annotation at
all.

**Note:** It is generally recommended not to document "optional" at all.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1 [(google.api.field_behavior) = REQUIRED];

  // The foreword for the book.
  string foreword = 2 [(google.api.field_behavior) = OPTIONAL];

  // The afterword for the book.
  string afterword = 3;  // Inconsistent use of OPTIONAL.
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1 [(google.api.field_behavior) = REQUIRED];

  // The foreword for the book.
  string foreword = 2;

  // The afterword for the book.
  string afterword = 3;
}
```

```proto
// Correct.
message Book {
  string name = 1 [(google.api.field_behavior) = REQUIRED];

  // The foreword for the book.
  string foreword = 2 [(google.api.field_behavior) = OPTIONAL];

  // The afterword for the book.
  string afterword = 3 [(google.api.field_behavior) = OPTIONAL];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0203::optional-consistency=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  string name = 1;

  // The foreword for the book.
  string foreword = 2 [(google.api.field_behavior) = OPTIONAL];

  // The afterword for the book.
  string afterword = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-203]: https://aip.dev/203
[aip.dev/not-precedent]: https://aip.dev/not-precedent
