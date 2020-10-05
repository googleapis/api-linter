---
rule:
  aip: 234
  name: [core, '0234', request-parent-reference]
  summary: |
    Batch Update requests should annotate the `parent` field with `google.api.resource_reference`.
permalink: /234/request-parent-reference
redirect_from:
  - /0234/request-parent-reference
---

# Batch Update methods: Parent resource reference

This rule enforces that all `BatchUpdate` requests have
`google.api.resource_reference` on their `string parent` field, as mandated in
[AIP-234][].

## Details

This rule looks at the `parent` field of any message matching `BatchUpdate*Request` and
complains if it does not have a `google.api.resource_reference` annotation.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message BatchUpdateBooksRequest {
  // The `google.api.resource_reference` annotation should also be included.
  string parent = 1;

  repeated UpdateBookRequest requests = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message BatchUpdateBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated UpdateBookRequest requests = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message BatchUpdateBooksRequest {
  // (-- api-linter: core::0234::request-parent-reference=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string parent = 1;

  repeated UpdateBookRequest requests = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-234]: https://aip.dev/234
[aip.dev/not-precedent]: https://aip.dev/not-precedent
