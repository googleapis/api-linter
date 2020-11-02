---
rule:
  aip: 165
  name: [core, '0165', request-filter-behavior]
  summary: |
    Purge requests should annotate the `filter` field with `google.api.field_behavior`.
permalink: /165/request-filter-behavior
redirect_from:
  - /0165/request-filter-behavior
---

# Purge requests: Filter field behavior

This rule enforces that all `Purge` requests have
`google.api.field_behavior` set to `REQUIRED` on their `string filter` field, as
mandated in [AIP-165][].

## Details

This rule looks at any message matching `Purge*Request` and complains if the
`filter` field does not have a `google.api.field_behavior` annotation with a
value of `REQUIRED`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message PurgeBooksRequest {
  string parent = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  // The `google.api.field_behavior` annotation should be included.
  string filter = 2;

  bool force = 3;
}
```

**Correct** code for this rule:

```proto
// Correct.
message PurgeBooksRequest {
  string parent = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  string filter = 2 [(google.api.field_behavior) = REQUIRED];

  bool force = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message PurgeBooksRequest {
  string parent = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  // (-- api-linter: core::0165::request-filter-behavior=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string filter = 2;

  bool force = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-165]: https://aip.dev/165
[aip.dev/not-precedent]: https://aip.dev/not-precedent
