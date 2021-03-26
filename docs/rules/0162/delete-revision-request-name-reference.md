---
rule:
  aip: 162
  name: [core, '0162', delete-revision-request-name-reference]
  summary: |
    Delete Revision requests should annotate the `name` field with `google.api.resource_reference`.
permalink: /162/delete-revision-request-name-reference
redirect_from:
  - /0162/delete-revision-request-name-reference
---

# Delete Revision requests: Name resource reference

This rule enforces that all Delete Revision requests have
`google.api.resource_reference` on their `string name` field, as mandated in
[AIP-162][].

## Details

This rule looks at the `name` field of any message matching `Delete*RevisionRequest`
and complains if it does not have a `google.api.resource_reference` annotation.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message DeleteBookRevisionRequest {
  // The `google.api.resource_reference` annotation should also be included.
  string name = 1 [(google.api.field_behavior) = REQUIRED];
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
  // (-- api-linter: core::0162::delete-revision-request-name-reference=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string name = 1 [(google.api.field_behavior) = REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent
