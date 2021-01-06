---
rule:
  aip: 162
  name: [core, '0162', tag-revision-request-tag-field]
  summary: Tag Revision RPCs must have a `tag` field in the request.
permalink: /162/tag-revision-request-tag-field
redirect_from:
  - /0162/tag-revision-request-tag-field
---

# Tag Revision requests: Tag field

This rule enforces that all Tag Revision methods have a `string tag`
field in the request message, as mandated in [AIP-162][].

## Details

This rule looks at any message matching `Tag*RevisionRequest` and complains if
either the `tag` field is missing or it has any type other than `string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should include a `string tag` field.
message TagBookRevisionRequest {
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

```proto
// Incorrect.
message TagBookRevisionRequest {
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  // Field type should be `string`.
  bytes tag = 2 [(google.api.field_behavior) = REQUIRED];
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

If you need to violate this rule, use a leading comment above the message (if
the `tag` field is missing) or above the field (if it is the wrong type).
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message TagBookRevisionRequest {
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  // (-- api-linter: core::0162::tag-revision-request-tag-field=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  bytes tag = 2 [(google.api.field_behavior) = REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent
