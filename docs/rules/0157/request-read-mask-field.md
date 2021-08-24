---
rule:
  aip: 157
  name: [core, '0157', request-read-mask-field]
  summary: Request read mask fields must have the correct type.
permalink: /157/request-read-mask-field
redirect_from:
  - /0157/request-read-mask-field
---

# Partial responses: Request read mask field

This rule enforces that all `read_mask` fields in requests have the correct
type, as mandated in [AIP-157][].

## Details

This rule looks at any message matching `*Request` that contains a `read_mask`
field, and complains if the field is not a singular `google.protobuf.FieldMask`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message GetBookRequest {
  string name = 1 [
    (google.api.resource_reference).type = "library.googleapis.com/Book",
    (google.api.field_behavior) = REQUIRED
  ];

  // Field type should be `google.protobuf.FieldMask`.
  string read_mask = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message GetBookRequest {
  string name = 1 [
    (google.api.resource_reference).type = "library.googleapis.com/Book",
    (google.api.field_behavior) = REQUIRED
  ];

  google.protobuf.FieldMask read_mask = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message GetBookRequest {
  string name = 1 [
    (google.api.resource_reference).type = "library.googleapis.com/Book",
    (google.api.field_behavior) = REQUIRED
  ];

  // (-- api-linter: core::0157::request-read-mask-field=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string read_mask = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-157]: https://aip.dev/157
[aip.dev/not-precedent]: https://aip.dev/not-precedent
