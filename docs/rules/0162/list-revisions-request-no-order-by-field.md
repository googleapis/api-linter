---
rule:
  aip: 162
  name: [core, '0162', list-revisions-request-no-order-by-field]
  summary: List Revisions requests should not have an `order_by` field.
permalink: /162/list-revisions-request-no-order-by-field
redirect_from:
  - /0162/list-revisions-request-no-order-by-field
---

# List Revisions requests: No `order_by` field

This rule enforces that List Revisions requests do not have a `string order_by`
field, as mandated in [AIP-162][].

## Details

This rule looks at any message matching `List*RevisionsRequest` and complains if
an `order_by` field is present, as revisions in the response must be ordered in
reverse chronological order.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message ListBookRevisionsRequest {
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  int32 page_size = 2;

  string page_token = 3;

  // Field should be removed.
  string order_by = 4;
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

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message ListBookRevisionsRequest {
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];

  int32 page_size = 2;

  string page_token = 3;

  // (-- api-linter: core::0162::list-revisions-request-no-order-by-field=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string order_by = 4;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent
