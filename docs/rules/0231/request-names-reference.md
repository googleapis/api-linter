---
rule:
  aip: 231
  name: [core, '0231', request-names-reference]
  summary: |
    Batch Get requests should annotate the `names` field with `google.api.resource_reference`.
permalink: /231/request-names-reference
redirect_from:
  - /0231/request-names-reference
---

# Batch Get methods: Resource reference

This rule enforces that all `BatchGet` requests have
`google.api.resource_reference` on their `repeated string names` field, as mandated in
[AIP-231][].

## Details

This rule looks at the `names` field of any message matching `BatchGet*Request` and
complains if it does not have a `google.api.resource_reference` annotation.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message BatchGetBooksRequest {
  // The `google.api.resource_reference` annotation should also be included.
  repeated string names = 1 [(google.api.field_behavior) = REQUIRED];
}
```

**Correct** code for this rule:

```proto
// Correct.
message BatchGetBooksRequest {
  repeated string names = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message BatchGetBooksRequest {
  // (-- api-linter: core::0231::request-names-reference=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  repeated string names = 1 [(google.api.field_behavior) = REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-231]: https://aip.dev/231
[aip.dev/not-precedent]: https://aip.dev/not-precedent
