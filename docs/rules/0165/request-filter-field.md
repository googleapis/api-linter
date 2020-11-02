---
rule:
  aip: 165
  name: [core, '0165', request-filter-field]
  summary: Purge RPCs must have a `filter` field in the request.
permalink: /165/request-filter-field
redirect_from:
  - /0165/request-filter-field
---

# Purge requests: Filter field

This rule enforces that all `Purge` methods have a `string filter`
field in the request message, as mandated in [AIP-165][].

## Details

This rule looks at any message matching `Purge*Request` and complains if
either the `filter` field is missing, or if it has any type other than `string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should include a `string filter` field.
message PurgeBooksRequest {
  string parent = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  bool force = 3;
}
```

```proto
// Incorrect.
message PurgeBooksRequest {
  string parent = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  // Field type should be `string`.
  bytes filter = 2 [(google.api.field_behavior) = REQUIRED];

  bool force = 3;
}
```

**Correct** code for this rule:

```proto
// Correct.
message PurgeBooksRequest {
  string parent = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  string filter = 2 [(google.api.field_behavior) = REQUIRED];

  bool force = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message (if
the `filter` field is missing) or above the field (if it is the wrong type).
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message PurgeBooksRequest {
  string parent = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  // (-- api-linter: core::0165::request-filter-field=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  bytes filter = 2 [(google.api.field_behavior) = REQUIRED];

  bool force = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-165]: https://aip.dev/165
[aip.dev/not-precedent]: https://aip.dev/not-precedent
