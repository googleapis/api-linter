---
rule:
  aip: 203
  name: [core, '0203', required-and-optional]
  summary: Required fields must not use the optional keyword.
permalink: /203/required-and-optional
redirect_from:
  - /0203/required-and-optional
---

# Required fields

This rule enforces that fields that are annotated as required do not use the
`optional` syntax, as mandated by [AIP-203][].

## Details

This rule looks for fields with a `google.api.field_behavior` annotation set to
`REQUIRED`, and complains if the field also uses the proto3 `optional` syntax.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;

  // Fields must not be optional and required.
  optional string title = 2 [(google.api.field_behavior) = REQUIRED];
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1;

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
  // (-- api-linter: core::0203::required-and-optional=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  optional string title = 2 [(google.api.field_behavior) = REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-203]: https://aip.dev/203
[aip.dev/not-precedent]: https://aip.dev/not-precedent
