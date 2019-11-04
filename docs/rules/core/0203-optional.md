---
rule:
  aip: 203
  name: [core, '0203', optional]
  summary: Optional fields may be annotated.
---

# Optional fields

This rule enforces that fields that are explicitly documented as optional also
have a machine-readable annotation, as mandated by [AIP-203][].

## Details

This rule looks at any field with "optional" (or similar forms) in the comment,
and complains if it does not have a `google.api.field_behavior` annotation.

**Note:** It is generally recommended not to document "optional" at all.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;

  // @Optional
  // The foreword for the book.
  string foreword = 2;
}
```

```proto
// Incorrect.
message Book {
  string name = 1;

  // Optional. The foreword for the book.
  string foreword = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1;

  // The foreword for the book.
  string foreword = 2 [(google.api.field_behavior) = OPTIONAL];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message Book {
  string name = 1;

  // Optional. The foreword for the book.
  // (-- api-linter: core::0203::optional
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string foreword = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-203]: https://aip.dev/203
[aip.dev/not-precedent]: https://aip.dev/not-precedent
