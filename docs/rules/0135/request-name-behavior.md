---
rule:
  aip: 135
  name: [core, '0135', request-name-behavior]
  summary: |
    Delete RPCs should annotate the `name` field with `google.api.field_behavior`.
permalink: /135/request-name-behavior
redirect_from:
  - /0135/request-name-behavior
---

# Delete methods: Field behavior

This rule enforces that all `Delete` standard methods have
`google.api.field_behavior` set to `REQUIRED` on their `string name` field, as
mandated in [AIP-135][].

## Details

This rule looks at any message matching `Delete*Request` and complains if the
`name` field does not have a `google.api.field_behavior` annotation with a
value of `REQUIRED`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message DeleteBookRequest {
  // The `google.api.field_behavior` annotation should also be included.
  string name = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Book"
  }];
}
```

**Correct** code for this rule:

```proto
// Correct.
message DeleteBookRequest {
  string name = 1 [
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
  // (-- api-linter: core::0135::request-name-behavior=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string name = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Book"
  }];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-135]: https://aip.dev/135
[aip.dev/not-precedent]: https://aip.dev/not-precedent
