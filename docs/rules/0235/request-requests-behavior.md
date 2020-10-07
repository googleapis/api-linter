---
rule:
  aip: 235
  name: [core, '0235', request-requests-behavior]
  summary: |
    Batch Delete requests should annotate the `requests` field with `google.api.field_behavior`.
permalink: /235/request-requests-behavior
redirect_from:
  - /0235/request-requests-behavior
---

# Batch Delete methods: `requests` field behavior

This rule enforces that all `BatchDelete` requests have
`google.api.field_behavior` set to `REQUIRED` on their `requests` field, as
mandated in [AIP-235][].

## Details

This rule looks at any message matching `BatchDelete*Request` and complains if the
`requests` field does not have a `google.api.field_behavior` annotation with a
value of `REQUIRED`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message BatchDeleteBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  // The `google.api.field_behavior` annotation should be included.
  repeated DeleteBookRequest requests = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message BatchDeleteBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated DeleteBookRequest requests = 2 [
    (google.api.field_behavior) = REQUIRED
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message BatchDeleteBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  // (-- api-linter: core::0235::request-requests-behavior=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  repeated DeleteBookRequest requests = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-235]: https://aip.dev/235
[aip.dev/not-precedent]: https://aip.dev/not-precedent
