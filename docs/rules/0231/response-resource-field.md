---
rule:
  aip: 231
  name: [core, '0231', response-resource-field]
  summary:
    Batch Get RPCs must have a repeated field for the resource in the response.
permalink: /231/response-resource-field
redirect_from:
  - /0231/response-resource-field
---

# Batch Get methods: Resource field

This rule enforces that all `BatchGet` methods have a repeated field in the
response message for the resource itself, as mandated in [AIP-231][].

## Details

This rule looks at any message matching `BatchGet*Response` and complains if
there is no repeated field of the resource's type with the expected field name.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message BatchGetBooksResponse {
  // `repeated Book books` is missing.
}
```

```proto
// Incorrect.
message BatchGetBooksResponse {
  Book books = 1;  // Field should be repeated.
}
```

**Correct** code for this rule:

```proto
// Correct.
message BatchGetBooksResponse {
  repeated Book books = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message (if
the resource field is missing) or above the field (if it is improperly named).
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0231::response-resource-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message BatchGetBooksResponse {
  Book books = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-231]: https://aip.dev/231
[aip.dev/not-precedent]: https://aip.dev/not-precedent
