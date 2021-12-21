---
rule:
  aip: 152
  name: [core, '0152', request-name-behavior]
  summary: |
    Run requests should annotate the `resource_name` field with `google.api.field_behavior`.
permalink: /152/request-name-behavior
redirect_from:
  - /0152/request-name-behavior
---

# Run requests: Field behavior

This rule enforces that all `Run` requests have
`google.api.field_behavior` set to `REQUIRED` on their `string resource_name` field, as
mandated in [AIP-152][].

## Details

This rule looks at any message matching `Run*JobRequest` and complains if the
`resource_name` field does not have a `google.api.field_behavior` annotation with a
value of `REQUIRED`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message RunWriteBookJobRequest {
  // The `google.api.field_behavior` annotation should also be included.
  string resource_name = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/WriteBookJob"
  }];
}
```

**Correct** code for this rule:

```proto
// Correct.
message RunWriteBookJobRequest {
  // The `google.api.field_behavior` annotation should also be included.
  string resource_name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference) = {
      type: "library.googleapis.com/WriteBookJob"
    }];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message RunWriteBookJobRequest {
  // (-- api-linter: core::0152::request-name-behavior=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string resource_name = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/WriteBookJob"
  }];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-152]: https://aip.dev/152
[aip.dev/not-precedent]: https://aip.dev/not-precedent
