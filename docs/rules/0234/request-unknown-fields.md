---
rule:
  aip: 234
  name: [core, '0234', request-unknown-fields]
  summary: Batch Update RPCs should not have unexpected fields in the request.
permalink: /234/request-unknown-fields
redirect_from:
  - /0234/request-unknown-fields
---

# Batch Update methods: Unknown request fields

This rule enforces that all `BatchUpdate` standard methods do not have unexpected
fields, as mandated in [AIP-234][].

## Details

This rule looks at any message matching `BatchUpdate*Request` and complains if it comes
across any fields other than:

- `bool allow_missing` ([AIP-134][])
- `string parent` ([AIP-234][])
- `string request_id` ([AIP-155][])
- `repeated Update*Request requests` ([AIP-234][])
- `google.protobuf.FieldMask update_mask` ([AIP-134][])
- `bool validate_only` ([AIP-163][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message BatchUpdateBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated UpdateBookRequest requests = 2 [
    (google.api.field_behavior) = REQUIRED
  ];

  string library_id = 3;  // Non-standard field.
}
```

**Correct** code for this rule:

```proto
// Correct.
message BatchUpdateBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated UpdateBookRequest requests = 2 [
    (google.api.field_behavior) = REQUIRED
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message BatchUpdateBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated UpdateBookRequest requests = 2 [
    (google.api.field_behavior) = REQUIRED
  ];

  // (-- api-linter: core::0234::request-unknown-fields=disabled
  //     aip.dev/not-precedent: We really need this field because reasons. --)
  string library_id = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-134]: https://aip.dev/134
[aip-155]: https://aip.dev/155
[aip-163]: https://aip.dev/163
[aip-234]: https://aip.dev/234
[aip.dev/not-precedent]: https://aip.dev/not-precedent
