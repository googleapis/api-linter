---
rule:
  aip: 162
  name: [core, '0162', delete-revision-request-name-behavior]
  summary: |
    Delete Revision requests should annotate the `name` field with `google.api.field_behavior`.
permalink: /162/delete-revision-request-name-behavior
redirect_from:
  - /0162/delete-revision-request-name-behavior
---

# Delete Revision requests: Name field behavior

This rule enforces that all Delete Revision requests have
`google.api.field_behavior` set to `REQUIRED` on their `string name` field, as
mandated in [AIP-162][].

## Details

This rule looks at any message matching `Delete*RevisionRequest` and complains if the
`name` field does not have a `google.api.field_behavior` annotation with a
value of `REQUIRED`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message DeleteBookRevisionRequest {
  // The `google.api.field_behavior` annotation should also be included.
  string name = 1 [
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

**Correct** code for this rule:

```proto
// Correct.
message DeleteBookRevisionRequest {
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
message DeleteBookRevisionRequest {
  // (-- api-linter: core::0162::delete-revision-request-name-behavior=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string name = 1 [
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent
