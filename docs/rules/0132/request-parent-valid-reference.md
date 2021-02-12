---
rule:
  aip: 132
  name: [core, '0132', request-parent-valid-reference]
  summary: |
    List RPCs should reference the parent resource, not the
    listed resource.
permalink: /132/request-parent-valid-reference
redirect_from:
  - /0132/request-parent-valid-reference
---

# List methods: Resource reference

This rule enforces that all `List` standard methods reference a resource other
than the resource being listed with the `google.api.resource_reference` on
their `string parent` field, as mandated in [AIP-132][].

## Details

This rule looks at the `parent` field of any message matching `List*Request`
and complains if the `google.api.resource_reference` annotation references
the resource being listed.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message ListBooksRequest {
  // The `google.api.resource_reference` should not reference the resource
  // being listed.
  string parent = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
  int32 page_size = 2;
  string page_token = 3;
}
```

**Correct** code for this rule:

```proto
// Correct.
message ListBooksRequest {
  string parent = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];
  int32 page_size = 2;
  string page_token = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0132::request-parent-valid-reference=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message ListBooksRequest {
  string parent = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-132]: https://aip.dev/132
[aip.dev/not-precedent]: https://aip.dev/not-precedent
