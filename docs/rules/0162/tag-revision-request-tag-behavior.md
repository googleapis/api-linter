---
rule:
  aip: 162
  name: [core, '0162', tag-revision-request-tag-behavior]
  summary: |
    Tag Revision requests should annotate the `tag` field with `google.api.field_behavior`.
permalink: /162/tag-revision-request-tag-behavior
redirect_from:
  - /0162/tag-revision-request-tag-behavior
---

# Tag Revision requests: Tag field behavior

This rule enforces that all Tag Revision requests have
`google.api.field_behavior` set to `REQUIRED` on their `string tag` field, as
mandated in [AIP-162][].

## Details

This rule looks at any message matching `Tag*RevisionRequest` and complains if the
`tag` field does not have a `google.api.field_behavior` annotation with a
value of `REQUIRED`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message TagBookRevisionRequest {
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  // The `google.api.field_behavior` annotation should be included.
  string tag = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message TagBookRevisionRequest {
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  string tag = 2 [(google.api.field_behavior) = REQUIRED];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message TagBookRevisionRequest {
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  // (-- api-linter: core::0162::tag-revision-request-tag-behavior=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string tag = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent
