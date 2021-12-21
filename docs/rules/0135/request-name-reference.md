---
rule:
  aip: 135
  name: [core, '0135', request-name-reference]
  summary: |
    Delete RPCs should annotate the `resource_name` field with `google.api.resource_reference`.
permalink: /135/request-name-reference
redirect_from:
  - /0135/request-name-reference
---

# Delete methods: Resource reference

This rule enforces that all `Delete` standard methods have
`google.api.resource_reference` on their `string resource_name` field, as mandated in
[AIP-135][].

## Details

This rule looks at the `resource_name` field of any message matching `Delete*Request`
and complains if it does not have a `google.api.resource_reference` annotation.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message DeleteBookRequest {
  // The `google.api.resource_reference` annotation should also be included.
  string resource_name = 1 [(google.api.field_behavior) = REQUIRED];
}
```

**Correct** code for this rule:

```proto
// Correct.
message DeleteBookRequest {
  string resource_name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message DeleteBookRequest {
  // (-- api-linter: core::0135::request-name-reference=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string resource_name = 1 [(google.api.field_behavior) = REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-135]: https://aip.dev/135
[aip.dev/not-precedent]: https://aip.dev/not-precedent
