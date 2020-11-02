---
rule:
  aip: 165
  name: [core, '0165', request-force-field]
  summary: Purge RPCs must have a `force` field in the request.
permalink: /165/request-force-field
redirect_from:
  - /0165/request-force-field
---

# Purge requests: Force field

This rule enforces that all `Purge` methods have a `bool force`
field in the request message, as mandated in [AIP-165][].

## Details

This rule looks at any message matching `Purge*Request` and complains if
either the `force` field is missing, or if it has any type other than `bool`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should include a `bool force` field.
message PurgeBooksRequest {
  string parent = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  string filter = 2 [(google.api.field_behavior) = REQUIRED];
}
```

```proto
// Incorrect.
message PurgeBooksRequest {
  string parent = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  string filter = 2 [(google.api.field_behavior) = REQUIRED];

  // Field type should be `bool`.
  int32 force = 3;
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
the `force` field is missing) or above the field (if it is the wrong type).
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message PurgeBooksRequest {
  string parent = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  string filter = 2 [(google.api.field_behavior) = REQUIRED];

  // (-- api-linter: core::0165::request-force-field=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  bool force = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-165]: https://aip.dev/165
[aip.dev/not-precedent]: https://aip.dev/not-precedent
