---
rule:
  aip: 213
  name: [core, '0213', common-types-messages]
  summary: Messages containing fields of a common type should use the common type.
permalink: /213/common-types-messages
redirect_from:
  - /0213/common-types-messages
---

# Common types messages

This rule enforces that messages that contain the fields of a common type use the
common type itself, as recommended in [AIP-213][].

## Details

This rule looks at the fields in a message, and complains if it finds a set of
field names that all belong to a common type.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;
  
  // Should use `google.type.Interval`.
  google.protobuf.Timestamp start_time = 2;
  google.protobuf.Timestamp end_time = 3;
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1;
  
  // Should use `google.type.Interval`.
  google.type.Interval interval = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the enum.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message Book {
  string name = 1;
  // (-- api-linter: core::0213::common-types-messages=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  google.protobuf.Timestamp expire_time = 2
    [(google.api.field_behavior) = OUTPUT_ONLY];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-213]: https://aip.dev/213
[aip.dev/not-precedent]: https://aip.dev/not-precedent
