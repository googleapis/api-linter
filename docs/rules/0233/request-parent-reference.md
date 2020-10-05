---
rule:
  aip: 233
  name: [core, '0233', request-parent-reference]
  summary: |
    Batch Create requests should annotate the `parent` field with `google.api.resource_reference`.
permalink: /233/request-parent-reference
redirect_from:
  - /0233/request-parent-reference
---

# Batch Create methods: Parent resource reference

This rule enforces that all `BatchCreate` requests have
`google.api.resource_reference` on their `string parent` field, as mandated in
[AIP-233][].

## Details

This rule looks at the `parent` field of any message matching `BatchCreate*Request` and
complains if it does not have a `google.api.resource_reference` annotation.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message BatchCreateBooksRequest {
  // The `google.api.resource_reference` annotation should also be included.
  string parent = 1;

  repeated CreateBookRequest requests = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message BatchCreateBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated CreateBookRequest requests = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message BatchCreateBooksRequest {
  // (-- api-linter: core::0233::request-parent-reference=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string parent = 1;

  repeated CreateBookRequest requests = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-233]: https://aip.dev/233
[aip.dev/not-precedent]: https://aip.dev/not-precedent
