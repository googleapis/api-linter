---
rule:
  aip: 164
  name: [core, '0164', request-name-reference]
  summary: |
    Undelete RPCs should annotate the `resource_name` field with `google.api.resource_reference`.
permalink: /164/request-name-reference
redirect_from:
  - /0164/request-name-reference
---

# Undelete methods: Resource reference

This rule enforces that all `Undelete` methods have
`google.api.resource_reference` on their `string resource_name` field, as mandated in
[AIP-164][].

## Details

This rule looks at the `resource_name` field of any message matching `Undelete*Request`
and complains if it does not have a `google.api.resource_reference` annotation.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message UndeleteBookRequest {
  // The `google.api.resource_reference` annotation should also be included.
  string resource_name = 1 [(google.api.field_behavior) = REQUIRED];
}
```

**Correct** code for this rule:

```proto
// Correct.
message UndeleteBookRequest {
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
message UndeleteBookRequest {
  // (-- api-linter: core::0164::request-name-reference=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string resource_name = 1 [(google.api.field_behavior) = REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-164]: https://aip.dev/164
[aip.dev/not-precedent]: https://aip.dev/not-precedent
