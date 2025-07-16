---
rule:
  aip: 142
  name: [core, '0142', duration-offset-comment]
  summary: Duration fields ending in `_offset` must have a clarifying comment.
permalink: /142/duration-offset-comment
redirect_from:
  - /0142/duration-offset-comment
---

# Duration offset comment

This rule enforces that `google.protobuf.Duration` fields ending in `_offset`
have a comment explaining what the offset is relative to, as mandated in
[AIP-142][].

## Details

This rule looks at each `google.protobuf.Duration` field that have a name
with an `_offset` suffix, and ensures that is has a leading
comment.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect: `start_offset` is a Duration with the `_offset` suffix
// but is missing the required comment.
message AudioSegment {
  google.protobuf.Duration start_offset = 1;

  // The total length of the segment.
  google.protobuf.Duration segment_duration = 2;
}
```

**Correct** code for this rule:

```proto
// Correct: `start_offset` is a Duration with the `_offset` suffix
// and has a clarifying comment.
message AudioSegment {
  // The duration relative to the start of the stream representing the
  // beginning of the segment.
  google.protobuf.Duration start_offset = 1;

  // The total length of the segment.
  google.protobuf.Duration segment_duration = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0142::duration-offset-comment=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message AudioSegment {
  google.protobuf.Duration start_offset = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-142]: https://aip.dev/142
[aip.dev/not-precedent]: https://aip.dev/not-precedent
