---
rule:
  aip: 135
  name: [core, '0135', request-unknown-fields]
  summary: Delete RPCs must not have unexpected required fields in the request.
permalink: /135/request-required-fields
redirect_from:
  - /0135/request-required-fields
---

# Delete methods: Required fields

This rule enforces that all `Delete` standard methods do not have unexpected
required fields, as mandated in [AIP-135][].

## Details

This rule looks at any message matching `Delete*Request` and complains if it
comes across any required fields other than:

- `string name` ([AIP-135][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message DeleteBookRequest {
  string name = 1 [(google.api.field_behavior) = REQUIRED];
  // Non-standard required field.
  bool allow_missing = 4 [(google.api.field_behavior) = REQUIRED];
}
```

**Correct** code for this rule:

```proto
// Correct.
message DeleteBookRequest {
  string name = 1 [(google.api.field_behavior) = REQUIRED];
  bool allow_missing = 4 [(google.api.field_behavior) = OPTIONAL];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message DeleteBookRequest {
  string name = 1 [(google.api.field_behavior) = REQUIRED];
  // (-- api-linter: core::0135::request-required-fields=disabled
  //     aip.dev/not-precedent: We really need this field to be required because
  // reasons. --)
  bool allow_missing = 4 [(google.api.field_behavior) = REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-135]: https://aip.dev/135
[aip.dev/not-precedent]: https://aip.dev/not-precedent
