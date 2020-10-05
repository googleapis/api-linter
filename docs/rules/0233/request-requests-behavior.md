---
rule:
  aip: 233
  name: [core, '0233', request-requests-behavior]
  summary: |
    Batch Create requests should annotate the `requests` field with `google.api.field_behavior`.
permalink: /233/request-requests-behavior
redirect_from:
  - /0233/request-requests-behavior
---

# Batch Create methods: `requests` field behavior

This rule enforces that all `BatchCreate` requests have
`google.api.field_behavior` set to `REQUIRED` on their `requests` field, as
mandated in [AIP-233][].

## Details

This rule looks at any message matching `BatchCreate*Request` and complains if the
`requests` field does not have a `google.api.field_behavior` annotation with a
value of `REQUIRED`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message BatchCreateBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  // The `google.api.field_behavior` annotation should be included.
  repeated CreateBookRequest requests = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message BatchCreateBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated CreateBookRequest requests = 2 [
    (google.api.field_behavior) = REQUIRED
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message BatchCreateBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  // (-- api-linter: core::0233::request-requests-behavior=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  repeated CreateBookRequest requests = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-233]: https://aip.dev/233
[aip.dev/not-precedent]: https://aip.dev/not-precedent
