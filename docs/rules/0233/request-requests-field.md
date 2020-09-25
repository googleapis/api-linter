---
rule:
  aip: 233
  name: [core, '0233', request-requests-field]
  summary: Batch Create RPCs should have a `requests` field in the request.
permalink: /233/request-requests-field
redirect_from:
  - /0233/request-requests-field
---

# Batch Create methods: Requests field

This rule enforces that all `BatchCreate` methods have a repeated `requests`
field, the type of which is the standard create request (like `Create*Request`),
in the request message, as mandated in [AIP-233][].

## Details

This rule looks at any message matching `BatchCreate*Request` and complains if
the `requests` field is missing, if it has any type other than `Create*Request`,
or if it is not `repeated`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message BatchCreateBooksRequest {
  string parent = 1;
  repeated CreateBookRequest req = 2;  // Field name should be `requests`.
}
```

```proto
// Incorrect.
message BatchCreateBooksRequest {
  string parent = 1;
  CreateBookRequest requests = 2;  // Field should be repeated.
}
```

**Correct** code for this rule:

```proto
// Correct.
message BatchCreateBooksRequest {
  string parent = 1;
  repeated CreateBookRequest requests = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message (if
the `requests` field is missing) or above the field (if it is the wrong type).
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0233::request-requests-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message BatchCreateBooksRequest {
  string parent = 1;
  repeated string books = 2; // should be "repeated CreateBookRequest requests"
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-233]: https://aip.dev/233
[aip.dev/not-precedent]: https://aip.dev/not-precedent
