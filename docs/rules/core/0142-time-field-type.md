---
rule:
  aip: 142
  name: [core, '0142', time-field-type]
  summary: Timestamps should use `google.protobuf.Timestamp`.
---

# Timestamp field type

This rule enforces that timestamps are represented with
`google.protobuf.Timestamp`, as mandated in [AIP-142][].

## Details

This rule looks at each field and looks for common suffixes that indicate that
the field may represent time, and indicates that the
`google.protobuf.Timestamp` type should be used instead.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;
  int32 publish_time_sec = 2;  // Should use `google.protobuf.Timestamp`.
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1;
  google.protobuf.Timestamp publish_time = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0142::time-field-type=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  string name = 1;
  int32 publish_time_sec = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-142]: https://aip.dev/142
[aip.dev/not-precedent]: https://aip.dev/not-precedent
