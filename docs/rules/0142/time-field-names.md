---
rule:
  aip: 142
  name: [core, '0142', time-field-names]
  summary: Timestamps should use `google.protobuf.Timestamp`.
permalink: /142/time-field-names
redirect_from:
  - /0142/time-field-names
---

# Timestamp field names

This rule enforces that timestamps are named using the imperative mood and with
a `_time` suffix, as mandated in [AIP-142][].

## Details

This rule looks at each `google.protobuf.Timestamp` field and ensures that it
has a `_time` suffix.

It also looks for common field names, regardless of type, and complains if they
are used. These are:

- created
- expired
- modified
- updated

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;
  google.protobuf.Timestamp published = 2;  // Should be `publish_time`.
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
// (-- api-linter: core::0142::time-field-names=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  string name = 1;
  google.protobuf.Timestamp published = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-142]: https://aip.dev/142
[aip.dev/not-precedent]: https://aip.dev/not-precedent
