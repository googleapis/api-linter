---
rule:
  aip: 234
  name: [core, '0234', request-requests-behavior]
  summary: |
    Batch Update requests should annotate the `requests` field with `google.api.field_behavior`.
permalink: /234/request-requests-behavior
redirect_from:
  - /0234/request-requests-behavior
---

# Batch Update methods: `requests` field behavior

This rule enforces that all `BatchUpdate` requests have
`google.api.field_behavior` set to `REQUIRED` on their `requests` field, as
mandated in [AIP-234][].

## Details

This rule looks at any message matching `BatchUpdate*Request` and complains if the
`requests` field does not have a `google.api.field_behavior` annotation with a
value of `REQUIRED`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message BatchUpdateBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  // The `google.api.field_behavior` annotation should be included.
  repeated UpdateBookRequest requests = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message BatchUpdateBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated UpdateBookRequest requests = 2 [
    (google.api.field_behavior) = REQUIRED
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message BatchUpdateBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  // (-- api-linter: core::0234::request-requests-behavior=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  repeated UpdateBookRequest requests = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-234]: https://aip.dev/234
[aip.dev/not-precedent]: https://aip.dev/not-precedent
