---
rule:
  aip: 133
  name: [core, '0133', request-required-fields-type]
  summary: Create RPCs must have correct types for required fields.
permalink: /133/request-required-fields-type
---

# Create methods: Required fields types

This rule enforces that all `Create` standard methods have the correct types for
required fields, as mandated in [AIP-133][].

## Details

This rule looks at any message matching `Create*Request` and complains if the
following required fields have the wrong type:

- `string parent` ([AIP-133][])
- `{Resource} {resource}` (must be a `message`)([AIP-133][])
- `string {resource}_id` ([AIP-133][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect: The `book` field should be a message type (Book).
message CreateBookRequest {
  string parent = 1 [(google.api.field_behavior) = REQUIRED];
  string book = 2 [(google.api.field_behavior) = REQUIRED];
  string book_id = 3;
}
```

**Correct** code for this rule:

```proto
// Correct.
message CreateBookRequest {
  string parent = 1 [(google.api.field_behavior) = REQUIRED];
  Book book = 2 [(google.api.field_behavior) = REQUIRED];
  string book_id = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message CreateBookRequest {
  // (-- api-linter: core::0133::request-required-fields-type=disabled
  //     aip.dev/not-precedent: We need to use bytes for the parent because
  // reasons. --)
  bytes parent = 1 [(google.api.field_behavior) = REQUIRED];
  Book book = 2 [(google.api.field_behavior) = REQUIRED];
  string book_id = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-133]: https://aip.dev/133
[aip.dev/not-precedent]: https://aip.dev/not-precedent
