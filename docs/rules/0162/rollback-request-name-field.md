---
rule:
  aip: 162
  name: [core, '0162', rollback-request-name-field]
  summary: Rollback RPCs must have a `resource_name` field in the request.
permalink: /162/rollback-request-name-field
redirect_from:
  - /0162/rollback-request-name-field
---

# Rollback requests: Resource Name field

This rule enforces that all `Rollback` methods have a `string resource_name`
field in the request message, as mandated in [AIP-162][].

## Details

This rule looks at any message matching `Rollback*Request` and complains if
either the `resource_name` field is missing or it has any type other than `string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should include a `string resource_name` field.
message RollbackBookRequest {
  string revision_id = 2 [(google.api.field_behavior) = REQUIRED];
}
```

```proto
// Incorrect.
message RollbackBookRequest {
  // Field type should be `string`.
  bytes resource_name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  string revision_id = 2 [(google.api.field_behavior) = REQUIRED];
}
```

**Correct** code for this rule:

```proto
// Correct.
message RollbackBookRequest {
  string resource_name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  string revision_id = 2 [(google.api.field_behavior) = REQUIRED];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message (if
the `resource_name` field is missing) or above the field (if it is the wrong type).
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message RollbackBookRequest {
  // (-- api-linter: core::0162::rollback-request-name-field=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  bytes resource_name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  string revision_id = 2 [(google.api.field_behavior) = REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent
