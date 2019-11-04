---
rule:
  aip: 203
  name: [core, '0203', output-only]
  summary: Output only fields should be annotated.
---

# Output only fields

This rule enforces that fields that are documented as output only also have a
machine-readable annotation, as mandated by [AIP-203][].

## Details

This rule looks at any field with "output only" (or similar forms) in the
comment, and complains if it does not have a `google.api.field_behavior`
annotation.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;

  // @OutputOnly
  // When the book was published.
  google.protobuf.Timestamp publish_time = 2;
}
```

```proto
// Incorrect.
message Book {
  string name = 1;

  // Output only. When the book was published.
  google.protobuf.Timestamp publish_time = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1;

  // When the book was published.
  google.protobuf.Timestamp publish_time = 2 [
    (google.api.field_behavior) = OUTPUT_ONLY
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message Book {
  string name = 1;

  // Immutable. The title of the book.
  // (-- api-linter: core::0203::output-only
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  // Output only. When the book was published.
  google.protobuf.Timestamp publish_time = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-203]: https://aip.dev/203
[aip.dev/not-precedent]: https://aip.dev/not-precedent
