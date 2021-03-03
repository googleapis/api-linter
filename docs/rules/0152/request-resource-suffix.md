---
rule:
  aip: 152
  name: [core, '0152', request-resource-suffix]
  summary: |
    Run requests should identify a resource type which ends in "Job".
permalink: /152/request-resource-suffix
redirect_from:
  - /0152/request-resource-suffix
---

# Run requests: Resource reference

This rule enforces that all `Run` methods specify a `type` that ends in "Job" 
in the `google.api.resource_reference` annotation of their `string name` field, 
as mandated in [AIP-152][].

## Details

This rule looks at the `name` field of any message matching `Run*JobRequest` 
and complains if the `type` in its `google.api.resource_reference` annotation
does not end in "Job".

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message RunWriteBookJobRequest {
  // The `type` of the `google.api.resource_reference` annotation should end in 
  // "Job".
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
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
  // (-- api-linter: core::0152::request-resource-suffix=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-152]: https://aip.dev/152
[aip.dev/not-precedent]: https://aip.dev/not-precedent
