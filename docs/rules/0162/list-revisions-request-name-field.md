---
rule:
  aip: 162
  name: [core, '0162', list-revisions-request-name-field]
  summary: List Revisions RPCs must have a `name` field in the request.
permalink: /162/list-revisions-request-name-field
redirect_from:
  - /0162/list-revisions-request-name-field
---

# List Revisions requests: Name field

This rule enforces that all List Revisions methods have a `string name`
field in the request message, as mandated in [AIP-162][].

## Details

This rule looks at any message matching `List*RevisionsRequest` and complains if
either the `name` field is missing or it has any type other than `string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should include a `string name` field.
message ListBookRevisionsRequest {
  int32 page_size = 1;

  string page_token = 2;
}
```

```proto
// Incorrect.
message ListBookRevisionsRequest {
  // Field type should be `string`.
  bytes name = 1 [
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
message ListBookRevisionsRequest {
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  int32 page_size = 2;

  string page_token = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message (if
the `name` field is missing) or above the field (if it is the wrong type).
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message ListBookRevisionsRequest {
  // (-- api-linter: core::0162::list-revisions-request-name-field=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  bytes name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  int32 page_size = 2;

  string page_token = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent
