---
rule:
  aip: 231
  name: [core, '0231', request-parent-field]
  summary: Batch Get RPCs must have a `parent` field in the request.
permalink: /231/request-parent-field
redirect_from:
  - /0231/request-parent-field
---

# Batch Get methods: Parent field

This rule enforces that all `BatchGet` methods have a `string parent` field in
the request message, as mandated in [AIP-231][].

## Details

This rule looks at any message matching `BatchGet*Request` and complains if
either the `parent` field is missing, or if it has any type other than
`string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message BatchGetBooksRequest {
  string publisher = 1;  // Field name should be `parent`.
  repeated string names = 2;
}
```

```proto
// Incorrect.
message BatchGetBookRequest {
  bytes parent = 1;  // Field type should be `string`.
  repeated string names = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message CreateBookRequest {
  string parent = 1;
  repeated string names = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message (if
the `parent` field is missing) or above the field (if it is the wrong type).
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0231::request-parent-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message BatchGetBooksRequest {
  string publisher = 1;  // Field name should be `parent`.
  repeated string names = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-231]: https://aip.dev/231
[aip.dev/not-precedent]: https://aip.dev/not-precedent
