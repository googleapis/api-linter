---
rule:
  aip: 134
  name: [core, '0134', request-unknown-fields]
  summary: Update RPCs must not have unexpected required fields in the request.
permalink: /134/request-required-fields
redirect_from:
  - /0134/request-required-fields
---

# Update methods: Required fields

This rule enforces that all `Update` standard methods do not have unexpected
required fields, as mandated in [AIP-134][].

## Details

This rule looks at any message matching `Update*Request` and complains if it
comes across any required fields other than:

- `{Resource} {resource}` ([AIP-134][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message UpdateBookRequest {
  Book book = 1 [(google.api.field_behavior) = REQUIRED];
  // Non-standard required field.
  bool allow_missing = 2 [(google.api.field_behavior) = REQUIRED];
}
```

**Correct** code for this rule:

```proto
// Correct.
message UpdateBookRequest {
  Book book = 1 [(google.api.field_behavior) = REQUIRED];
  bool allow_missing = 2 [(google.api.field_behavior) = OPTIONAL];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message UpdateBookRequest {
  Book book = 1 [(google.api.field_behavior) = REQUIRED];
  // (-- api-linter: core::0134::request-required-fields=disabled
  //     aip.dev/not-precedent: We really need this field to be required because
  // reasons. --)
  bool allow_missing = 2 [(google.api.field_behavior) = REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-134]: https://aip.dev/134
[aip.dev/not-precedent]: https://aip.dev/not-precedent
