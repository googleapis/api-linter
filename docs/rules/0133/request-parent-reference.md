---
rule:
  aip: 133
  name: [core, '0133', request-parent-reference]
  summary: |
    Create RPCs should annotate the `parent` field with `google.api.resource_reference`.
permalink: /133/request-parent-reference
redirect_from:
  - /0133/request-parent-reference
---

# Create methods: Resource reference

This rule enforces that all `Create` standard methods have
`google.api.resource_reference` on their `string parent` field, as mandated in
[AIP-133][].

## Details

This rule looks at the `parent` field of any message matching `Create*Request`
and complains if it does not have a `google.api.resource_reference` annotation.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message CreateBookRequest {
  // The `google.api.resource_reference` annotation should also be included.
  string parent = 1 [(google.api.field_behavior) = REQUIRED];
  Book book = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message CreateBookRequest {
  string parent = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Publisher"
  ];
  Book book = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message CreateBookRequest {
  // (-- api-linter: core::0133::request-parent-reference=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string parent = 1 [(google.api.field_behavior) = REQUIRED];
  Book book = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-133]: https://aip.dev/133
[aip.dev/not-precedent]: https://aip.dev/not-precedent
