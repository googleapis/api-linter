---
rule:
  aip: 165
  name: [core, '0165', response-purge-sample-field]
  summary: Purge RPCs must have a `purge_sample` field in the response.
permalink: /165/response-purge-sample-field
redirect_from:
  - /0165/response-purge-sample-field
---

# Purge responses: Purge sample field

This rule enforces that all `Purge` methods have a `repeated string
purge_sample` field in the response message, as mandated in [AIP-165][].

## Details

This rule looks at any message matching `Purge*Response` and complains if
either the `purge_sample` field is missing, or if it has any type other than
`repeated string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should include a `repeated string purge_sample` field.
message PurgeBooksResponse {
  int32 purge_count = 1;
}
```

```proto
// Incorrect.
message PurgeBooksResponse {
  int64 purge_count = 1;

  // Field type should be `repeated string`.
  repeated bytes purge_sample = 2 [
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
the `purge_sample` field is missing) or above the field (if it is the wrong
type). Remember to also include an [aip.dev/not-precedent][] comment explaining
why.

```proto
message PurgeBooksResponse {
  int32 purge_count = 1;

  // (-- api-linter: core::0165::response-purge-sample-field=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  repeated bytes purge_sample = 2 [
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-165]: https://aip.dev/165
[aip.dev/not-precedent]: https://aip.dev/not-precedent
