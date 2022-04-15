---
rule:
  aip: 133
  name: [core, '0133', resource-reference-type]
  summary: Create should use a `child_type` reference to the created resource.
permalink: /133/resource-reference-type
redirect_from:
  - /0133/resource-reference-type
---

# Create methods: Parent field resource reference

This rule enforces that all `Create` standard methods with a `string parent`
field use a proper `google.api.resource_reference`, that being either a
`child_type` referring to the created resource or a `type` referring directly
to the parent resource, as mandated in [AIP-133][].

## Details

This rule looks at any message matching `Create*Request` and complains if the 
`google.api.resource_reference` on the `parent` field refers to the wrong
resource.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message CreateBooksRequest {
  // `child_type` should be used instead of `type` when referring to the
  // created resource on a parent field.
  string parent = 1 [(google.api.resource_reference).type = "library.googleapis.com/Book"];
  Book book = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message CreateBooksRequest {
  string parent = 1 [(google.api.resource_reference).child_type = "library.googleapis.com/Book"];
  Book book = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message CreateBooksRequest {
  // (-- api-linter: core::0133::resource-reference-type=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string parent = 1 [(google.api.resource_reference).type = "library.googleapis.com/Book"];
  Book book = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-133]: https://aip.dev/133
[aip.dev/not-precedent]: https://aip.dev/not-precedent
