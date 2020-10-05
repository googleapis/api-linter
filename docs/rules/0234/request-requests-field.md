---
rule:
  aip: 234
  name: [core, '0234', request-requests-field]
  summary: Batch Update RPCs should have a `requests` field in the request.
permalink: /234/request-requests-field
redirect_from:
  - /0234/request-requests-field
---

# Batch Update methods: Requests field

This rule enforces that all `BatchUpdate` methods have a repeated `requests`
field, the type of which is the standard Update request (`Update*Request`)
in the request message, as mandated in [AIP-234][].

## Details

This rule looks at any message matching `BatchUpdate*Request` and complains if
the `requests` field is missing, if it has any type other than `Update*Request`,
or if it is not `repeated`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message BatchUpdateBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated UpdateBookRequest req = 2;  // Field name should be `requests`.
}
```

```proto
// Incorrect.
message BatchUpdateBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  UpdateBookRequest requests = 2;  // Field should be repeated.
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
the `requests` field is missing) or above the field (if it is the wrong type).
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0234::request-requests-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message BatchUpdateBooksRequest {
  string parent = 1 [
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  repeated string books = 2; // should be "repeated UpdateBookRequest requests"
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-234]: https://aip.dev/234
[aip.dev/not-precedent]: https://aip.dev/not-precedent
