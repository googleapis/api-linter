---
rule:
  aip: 132
  name: [core, '0132', request-parent-behavior]
  summary: |
    List RPCs should annotate the `parent` field with `google.api.field_behavior`.
permalink: /132/request-parent-behavior
redirect_from:
  - /0132/request-parent-behavior
---

# List methods: Field behavior

This rule enforces that all `List` standard methods have
`google.api.field_behavior` set to `REQUIRED` on their `string parent` field,
as mandated in [AIP-132][].

## Details

This rule looks at any message matching `List*Request` and complains if the
`parent` field does not have a `google.api.field_behavior` annotation with a
value of `REQUIRED`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message ListBooksRequest {
  // The `google.api.field_behavior` annotation should also be included.
  string parent = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Publisher"
  }];
  int32 page_size = 2;
  string page_token = 3;
}
```

**Correct** code for this rule:

```proto
// Correct.
message ListBooksRequest {
  string parent = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Publisher"
  ];
  int32 page_size = 2;
  string page_token = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0132::request-parent-behavior=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message ListBooksRequest {
  string parent = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Publisher"
  }];
  int32 page_size = 2;
  string page_token = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-132]: https://aip.dev/132
[aip.dev/not-precedent]: https://aip.dev/not-precedent
