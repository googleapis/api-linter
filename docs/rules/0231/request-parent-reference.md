---
rule:
  aip: 231
  name: [core, '0231', request-parent-reference]
  summary: |
    Batch Get requests should annotate the `parent` field with `google.api.resource_reference`.
permalink: /231/request-parent-reference
redirect_from:
  - /0231/request-parent-reference
---

# Batch Get methods: Parent resource reference

This rule enforces that all `BatchGet` requests have
`google.api.resource_reference` on their `string parent` field, as mandated in
[AIP-231][].

## Details

This rule looks at the `parent` field of any message matching `BatchGet*Request` and
complains if it does not have a `google.api.resource_reference` annotation.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message BatchGetBooksRequest {
  // The `google.api.resource_reference` annotation should also be included.
  string parent = 1;

  repeated string names = 2 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

**Correct** code for this rule:

```proto
// Correct.
message BatchGetBooksRequest {
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

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message BatchGetBooksRequest {
  // (-- api-linter: core::0231::request-parent-reference=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string parent = 1;

  repeated string names = 2 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-231]: https://aip.dev/231
[aip.dev/not-precedent]: https://aip.dev/not-precedent
