---
rule:
  aip: 133
  name: [core, '0133', request-required-fields]
  summary: Create RPCs must not have unexpected required fields in the request.
permalink: /133/request-required-fields
redirect_from:
  - /0133/request-required-fields
---

# Create methods: Required fields

This rule enforces that all `Create` standard methods do not have unexpected
required fields, as mandated in [AIP-133][].

## Details

This rule looks at any message matching `Create*Request` and complains if it
comes across any required fields other than:

- `string parent` ([AIP-133][])
- `{Resource} {resource}` ([AIP-133][])
- `string {resource}_id` ([AIP-133][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message CreateBookRequest {
  string parent = 1;
  Book book = 2;
  string book_id = 3;
  // Non-standard required field.
  string validate_only = 4 [(google.api.field_behavior) = REQUIRED];
}
```

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
  string parent = 1;
  Book book = 2;
  string book_id = 3;
  string validate_only = 4 [(google.api.field_behavior) = OPTIONAL];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message CreateBookRequest {
  string parent = 1;
  Book book = 2;
  string book_id = 3;

  // (-- api-linter: core::0133::request-required-fields=disabled
  //     aip.dev/not-precedent: We really need this field to be required because
  // reasons. --)
  string validate_only = 4 [(google.api.field_behavior) = REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-133]: https://aip.dev/133
[aip.dev/not-precedent]: https://aip.dev/not-precedent
