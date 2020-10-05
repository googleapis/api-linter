---
rule:
  aip: 231
  name: [core, '0231', request-names-field]
  summary: Batch Get RPCs should have a `names` field in the request.
permalink: /231/request-names-field
redirect_from:
  - /0231/request-names-field
---

# Batch Get methods: Names field

This rule enforces that all `BatchGet` methods have a `repeated string names`
field in the request message, as mandated in [AIP-231][].

## Details

This rule looks at any message matching `BatchGet*Request` and complains if
either the `names` field is missing, if it has any type other than `string`, or
if it is not `repeated`.

Alternatively, if there is a `repeated GetBookRequest requests` field, this is
accepted in its place.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message BatchGetBooksRequest {
  string parent = 1;
  repeated string books = 2;  // Field name should be `names`.
}
```

```proto
// Incorrect.
message BatchGetBooksRequest {
  string parent = 1;
  string names = 2;  // Field should be repeated.
}
```

**Correct** code for this rule:

```proto
// Correct.
message BatchGetBooksRequest {
  string parent = 1;
  repeated string names = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message (if
the `names` field is missing) or above the field (if it is the wrong type).
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0231::request-names-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message BatchGetBooksRequest {
  string parent = 1;
  repeated string books = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-231]: https://aip.dev/231
[aip.dev/not-precedent]: https://aip.dev/not-precedent
