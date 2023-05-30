---
rule:
  aip: 131
  name: [core, '0131', request-required-fields]
  summary: Get RPCs must not have unexpected required fields in the request.
permalink: /131/request-required-fields
redirect_from:
  - /0131/request-required-fields
---

# Get methods: Required fields

This rule enforces that all `Get` standard methods do not have unexpected
required fields, as mandated in [AIP-131][].

## Details

This rule looks at any message matching `Get*Request` and complains if it
comes across any required fields other than:

- `string name` ([AIP-131][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message GetBookRequest {
  // The name of the book to retrieve.
  // Format: publishers/{publisher}/books/{book}
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference) = {
      type: "library.googleapis.com/Book"
  }];

  // Non-standard required field.
   google.protobuf.FieldMask read_mask = 2 [(google.api.field_behavior) = REQUIRED];
}
```

**Correct** code for this rule:

```proto
// Correct.
message GetBookRequest {
  // The name of the book to retrieve.
  // Format: publishers/{publisher}/books/{book}
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference) = {
      type: "library.googleapis.com/Book"
  }];

  google.protobuf.FieldMask read_mask = 2 [(google.api.field_behavior) = OPTIONAL];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message GetBookRequest {
  // The name of the book to retrieve.
  // Format: publishers/{publisher}/books/{book}
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference) = {
      type: "library.googleapis.com/Book"
  }];

  // (-- api-linter: core::0131::request-required-fields=disabled
  //     aip.dev/not-precedent: We really need this field to be required because
  //     reasons. --)
   google.protobuf.FieldMask read_mask = 2 [(google.api.field_behavior) = REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-131]: https://aip.dev/131
[aip.dev/not-precedent]: https://aip.dev/not-precedent
