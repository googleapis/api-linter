---
rule:
  aip: 234
  name: [core, '0234', request-parent-field]
  summary: Batch Update RPCs must have a `parent` field in the request.
permalink: /234/request-parent-field
redirect_from:
  - /0234/request-parent-field
---

# Batch Update methods: Parent field

This rule enforces that all `BatchUpdate` methods have a `string parent` field
in the request message, as mandated in [AIP-234][].

## Details

This rule looks at any message matching `BatchUpdate*Request` and complains if
either the `parent` field is missing, or if it has any type other than
`string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message BatchUpdateBooksRequest {
  string publisher = 1;  // Field name should be `parent`.
  repeated UpdateBookRequest requests = 2;
}
```

```proto
// Incorrect.
message BatchUpdateBooksRequest {
  bytes parent = 1;  // Field type should be `string`.
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

If you need to violate this rule, use a leading comment above the message (if
the `parent` field is missing) or above the field (if it is the wrong type).
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0234::request-parent-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message BatchUpdateBooksRequest {
  string publisher = 1;  // Field name should be `parent`.
  repeated UpdateBookRequest requests = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-234]: https://aip.dev/234
[aip.dev/not-precedent]: https://aip.dev/not-precedent
