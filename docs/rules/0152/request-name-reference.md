---
rule:
  aip: 152
  name: [core, '0152', request-name-reference]
  summary: |
    Run requests should annotate the `name` field with `google.api.resource_reference`.
permalink: /152/request-name-reference
redirect_from:
  - /0152/request-name-reference
---

# Run requests: Resource reference

This rule enforces that all `RUn` methods have
`google.api.resource_reference` on their `string name` field, as mandated in
[AIP-152][].

## Details

This rule looks at the `name` field of any message matching `Run*JobRequest`
and complains if it does not have a `google.api.resource_reference` annotation.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message RunWriteBookJobRequest {
  // The `google.api.resource_reference` annotation should also be included.
  string name = 1 [(google.api.field_behavior) = REQUIRED];
}
```

**Correct** code for this rule:

```proto
// Correct.
message RunWriteBookJobRequest {
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/WriteBookJob"
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message RunWriteBookJobRequest {
  // (-- api-linter: core::0152::request-name-reference=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string name = 1 [(google.api.field_behavior) = REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-152]: https://aip.dev/152
[aip.dev/not-precedent]: https://aip.dev/not-precedent
