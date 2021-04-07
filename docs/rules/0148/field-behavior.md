---
rule:
  aip: 148
  name: [core, '0148',field-behavior]
  summary: Standard resource fields should have the correct field behavior.
permalink: /148/field-behavior
redirect_from:
  - /0148/field-behavior
---

# Standard resource fields: Field behavior

This rule enforces that all standard resource fields have the correct
`google.api.field_behavior`, as mandated in [AIP-148][].

## Details

This rule looks at any message with a `google.api.resource` annotation, and
complains if any of the following fields does not have a
`google.api.field_behavior` annotation with a value of `OUTPUT_ONLY`:

- `create_time`
- `delete_time`
- `uid`
- `update_time`

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  string name = 1;

  // The `google.api.field_behavior` annotation should be `OUTPUT_ONLY`.
  google.protobuf.Timestamp create_time = 2;

  // The `google.api.field_behavior` annotation should be `OUTPUT_ONLY`.
  google.protobuf.Timestamp update_time = 3;

  // The `google.api.field_behavior` annotation should be `OUTPUT_ONLY`.
  google.protobuf.Timestamp delete_time = 4;

  // The `google.api.field_behavior` annotation should be `OUTPUT_ONLY`.
  string uid = 5;
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  string name = 1;

  google.protobuf.Timestamp create_time = 2 [(google.api.field_behavior) = OUTPUT_ONLY];

  google.protobuf.Timestamp update_time = 3 [(google.api.field_behavior) = OUTPUT_ONLY];

  google.protobuf.Timestamp delete_time = 4 [(google.api.field_behavior) = OUTPUT_ONLY];

  string uid = 5 [(google.api.field_behavior) = OUTPUT_ONLY];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  string name = 1;

  // (-- api-linter: core::0148::field-behavior=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  google.protobuf.Timestamp create_time = 2;

  // (-- api-linter: core::0148::field-behavior=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  google.protobuf.Timestamp update_time = 3;

  // (-- api-linter: core::0148::field-behavior=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  google.protobuf.Timestamp delete_time = 4;

  // (-- api-linter: core::0148::field-behavior=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string uid = 5;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-148]: https://aip.dev/148
[aip.dev/not-precedent]: https://aip.dev/not-precedent
