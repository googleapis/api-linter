---
rule:
  aip: 131
  name: [core, '0131', request-name-reference-type]
  summary: |
    The `google.api.resource_reference` on the `resource_name` field of a Get RPC request
    message should use `type`, not `child_type`.
permalink: /131/request-name-reference-type
redirect_from:
  - /0131/request-name-reference-type
---

# Get methods: Resource reference

This rule enforces that the `google.api.resource_reference` on the `resource_name` field
of a Get RPC request message uses `type`, not `child_type`, as suggested in
[AIP-131][].

## Details

This rule looks at the `google.api.resource_reference` annotation on the `resource_name`
field of any message matching `Get*Request` and complains if it does not use a
direct `type` reference.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message GetBookRequest {
  // The `google.api.resource_reference` annotation should be a direct `type`
  // reference.
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];
}
```

**Correct** code for this rule:

```proto
// Correct.
message GetBookRequest {
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message GetBookRequest {
  // (-- api-linter: core::0131::request-name-reference-type=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-131]: https://aip.dev/131
[aip.dev/not-precedent]: https://aip.dev/not-precedent
