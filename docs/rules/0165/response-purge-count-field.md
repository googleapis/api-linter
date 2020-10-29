---
rule:
  aip: 165
  name: [core, '0165', response-purge-count-field]
  summary: Purge RPCs must have a `purge_count` field in the response.
permalink: /165/response-purge-count-field
redirect_from:
  - /0165/response-purge-count-field
---

# Purge responses: Purge count field

This rule enforces that all `Purge` methods have an `int32 purge_count`
field in the response message, as mandated in [AIP-165][].

## Details

This rule looks at any message matching `Purge*Response` and complains if
either the `purge_count` field is missing, or if it has any type other than
`int32`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should include an `int32 purge_count` field.
message PurgeBooksResponse {
  repeated string purge_sample = 2 [
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

```proto
// Incorrect.
message PurgeBooksResponse {
  // Field type should be `int32`.
  int64 purge_count = 1;

  repeated string purge_sample = 2 [
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
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

If you need to violate this rule, use a leading comment above the message (if
the `purge_count` field is missing) or above the field (if it is the wrong
type). Remember to also include an [aip.dev/not-precedent][] comment explaining
why.

```proto
message PurgeBooksResponse {
  // (-- api-linter: core::0165::response-purge-count-field=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  int64 purge_count = 1;

  repeated string purge_sample = 2 [
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-165]: https://aip.dev/165
[aip.dev/not-precedent]: https://aip.dev/not-precedent
