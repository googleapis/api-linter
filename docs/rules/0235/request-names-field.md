---
rule:
  aip: 235
  name: [core, '0235', request-names-field]
  summary: Batch Delete RPCs should have a `names` field in the request.
permalink: /235/request-names-field
redirect_from:
  - /0235/request-names-field
---

# Batch Delete methods: Names field

This rule enforces that all `BatchDelete` methods have a `repeated string names`
field in the request message, as mandated in [AIP-235][].

## Details

This rule looks at any message matching `BatchDelete*Request` and complains if
the `names` field is missing, if it has any type other than `string`, or
if it is not `repeated`.

Alternatively, if there is a `repeated DeleteBookRequest requests` field, this is
accepted in its place.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message BatchDeleteBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated string books = 2;  // Field name should be `names`.
}
```

```proto
// Incorrect.
message BatchDeleteBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  string names = 2;  // Field should be repeated.
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
the `names` field is missing) or above the field (if it is the wrong type).
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0235::request-names-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message BatchDeleteBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated string books = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-235]: https://aip.dev/235
[aip.dev/not-precedent]: https://aip.dev/not-precedent
