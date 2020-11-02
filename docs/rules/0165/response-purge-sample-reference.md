---
rule:
  aip: 165
  name: [core, '0165', response-purge-sample-reference]
  summary: |
    Purge responses should annotate the `purge_sample` field with `google.api.resource_reference`.
permalink: /165/response-purge-sample-reference
redirect_from:
  - /0165/response-purge-sample-reference
---

# Purge responses: Purge sample resource reference

This rule enforces that all `Purge` responses have
`google.api.resource_reference` on their `repeated string purge_sample` field,
as mandated in [AIP-165][].

## Details

This rule looks at the `purge_sample` field of any message matching
`Purge*Response` and complains if it does not have a
`google.api.resource_reference` annotation.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message PurgeBooksResponse {
  int32 purge_count = 1;

  // The `google.api.resource_reference` annotation should be included.
  repeated string purge_sample = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message PurgeBooksResponse {
  int32 purge_count = 1;

  repeated string purge_sample = 2 [
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message PurgeBooksResponse {
  int32 purge_count = 1;

  // (-- api-linter: core::0165::response-purge-sample-reference=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  repeated string purge_sample = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-165]: https://aip.dev/165
[aip.dev/not-precedent]: https://aip.dev/not-precedent
