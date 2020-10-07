---
rule:
  aip: 235
  name: [core, '0235', request-parent-field]
  summary: Batch Delete RPCs must have a `parent` field in the request.
permalink: /235/request-parent-field
redirect_from:
  - /0235/request-parent-field
---

# Batch Delete methods: Parent field

This rule enforces that all `BatchDelete` methods have a `string parent` field in
the request message, as mandated in [AIP-235][].

## Details

This rule looks at any message matching `BatchDelete*Request` and complains if
either the `parent` field is missing, or if it has any type other than
`string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message BatchDeleteBooksRequest {
  string publisher = 1;  // Field name should be `parent`.

  repeated string names = 2 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

```proto
// Incorrect.
message BatchDeleteBooksRequest {
  bytes parent = 1;  // Field type should be `string`.

  repeated string names = 2 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

**Correct** code for this rule:

```proto
// Correct.
message BatchDeleteBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated string names = 2 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message (if
the `parent` field is missing) or above the field (if it is the wrong type).
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0235::request-parent-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message BatchDeleteBooksRequest {
  string publisher = 1;  // Field name should be `parent`.

  repeated string names = 2 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-235]: https://aip.dev/235
[aip.dev/not-precedent]: https://aip.dev/not-precedent
