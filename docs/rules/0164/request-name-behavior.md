---
rule:
  aip: 164
  name: [core, '0164', request-name-behavior]
  summary: |
    Undelete RPCs should annotate the `resource_name` field with `google.api.field_behavior`.
permalink: /164/request-name-behavior
redirect_from:
  - /0164/request-name-behavior
---

# Undelete methods: Field behavior

This rule enforces that all `Undelete` methods have
`google.api.field_behavior` set to `REQUIRED` on their `string resource_name` field, as
mandated in [AIP-164][].

## Details

This rule looks at any message matching `Undelete*Request` and complains if the
`resource_name` field does not have a `google.api.field_behavior` annotation with a
value of `REQUIRED`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message UndeleteBookRequest {
  // The `google.api.field_behavior` annotation should also be included.
  string resource_name = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Book"
  }];
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
  // (-- api-linter: core::0164::request-name-behavior=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string resource_name = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Book"
  }];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-164]: https://aip.dev/164
[aip.dev/not-precedent]: https://aip.dev/not-precedent
