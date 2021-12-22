---
rule:
  aip: 162
  name: [core, '0162', commit-request-name-reference]
  summary: |
    Commit requests should annotate the `resource_name` field with `google.api.resource_reference`.
permalink: /162/commit-request-name-reference
redirect_from:
  - /0162/commit-request-name-reference
---

# Commit requests: Resource Name resource reference

This rule enforces that all `Commit` requests have
`google.api.resource_reference` on their `string resource_name` field, as mandated in
[AIP-162][].

## Details

This rule looks at the `resource_name` field of any message matching `Commit*Request`
and complains if it does not have a `google.api.resource_reference` annotation.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message CommitBookRequest {
  // The `google.api.resource_reference` annotation should also be included.
  string resource_name = 1 [(google.api.field_behavior) = REQUIRED];
}
```

**Correct** code for this rule:

```proto
// Correct.
message CommitBookRequest {
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
message CommitBookRequest {
  // (-- api-linter: core::0162::commit-request-name-reference=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string resource_name = 1 [(google.api.field_behavior) = REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-162]: https://aip.dev/162
[aip.dev/not-precedent]: https://aip.dev/not-precedent
