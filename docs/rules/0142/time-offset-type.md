---
rule:
  aip: 142
  name: [core, '0142', time-offset-type]
  summary: Fields ending in `_time_offset` must be of type `google.protobuf.Duration`.
permalink: /142/time-offset-type
redirect_from:
  - /0142/time-offset-type
---

# Time offset type

This rule enforces that fields with names ending in `_time_offset` use the
`google.protobuf.Duration` type, as mandated in [AIP-142][].

## Details
This rule looks at each field with a name ending in `_time_offset` and
verifies that its type is `google.protobuf.Duration`.


## Examples

**Incorrect** code for this rule:

```proto
// Incorrect: `start_time_offset` ends in `_time_offset` but is not a `google.protobuf.Duration`.
message AudioSegment {
  // The duration relative to the start of the stream representing the
  // beginning of the segment.
  int64 start_time_offset = 1;
}
```

**Correct** code for this rule:

```proto
// Correct: `start_time_offset` is a `google.protobuf.Duration`.
message AudioSegment {
  // The duration relative to the start of the stream representing the
  // beginning of the segment.
  google.protobuf.Duration start_time_offset = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0142::time-offset-type=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message AudioSegment {
  // The duration relative to the start of the stream representing the
  // beginning of the segment.
  int64 start_time_offset = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-142]: https://aip.dev/142
[aip.dev/not-precedent]: https://aip.dev/not-precedent
