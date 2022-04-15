---
rule:
  aip: 132
  name: [core, '0132', resource-reference-type]
  summary: List should use a `child_type` reference to the paginated resource.
permalink: /132/resource-reference-type
redirect_from:
  - /0132/resource-reference-type
---

# List methods: Parent field resource reference

This rule enforces that all `List` standard methods with a `string parent`
field use a proper `google.api.resource_reference`, that being either a
`child_type` referring to the pagianted resource or a `type` referring directly
to the parent resource, as mandated in [AIP-132][].

## Details

This rule looks at any message matching `List*Request` and complains if the 
`google.api.resource_reference` on the `parent` field refers to the wrong
resource.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message ListBooksRequest {
  // `child_type` should be used instead of `type` when referring to the
  // paginated resource on a parent field.
  string parent = 1 [(google.api.resource_reference).type = "library.googleapis.com/Book"];
  int32 page_size = 2;
  string page_token = 3;
}
```

**Correct** code for this rule:

```proto
// Correct.
message ListBooksRequest {
  string parent = 1 [(google.api.resource_reference).child_type = "library.googleapis.com/Book"];
  int32 page_size = 2;
  string page_token = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message ListBooksRequest {
  // (-- api-linter: core::0132::resource-reference-type=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string parent = 1 [(google.api.resource_reference).type = "library.googleapis.com/Book"];
  int32 page_size = 2;
  string page_token = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-132]: https://aip.dev/132
[aip.dev/not-precedent]: https://aip.dev/not-precedent
