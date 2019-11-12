---
rule:
  aip: 203
  name: [core, '0203', optional-conflict]
  summary: Optional fields must not have any other field behavior.
---

# Optional fields: Conflicts

This rule enforces that fields that are annotated as `OPTIONAL` do not have any
other annotation (such as `REQUIRED`), as mandated by [AIP-203][].

## Details

This rule looks at any field with a `google.api.field_behavior` annotation of
`OPTIONAL`, and complains if it also finds any other field behavior.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;

  // The foreword for the book.
  string foreword = 2 [
    (google.api.field_behavior) = OPTIONAL,  // "Optional" can not co-exist with other field behaviors.
    (google.api.field_behavior) = IMMUTABLE];
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
  // (-- api-linter: core::0203::optional-conflict
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string foreword = 2 [
    (google.api.field_behavior) = OPTIONAL,
    (google.api.field_behavior) = IMMUTABLE];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-203]: https://aip.dev/203
[aip.dev/not-precedent]: https://aip.dev/not-precedent
