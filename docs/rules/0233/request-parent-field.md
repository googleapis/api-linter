---
rule:
  aip: 233
  name: [core, '0233', request-parent-field]
  summary: Batch Create RPCs must have a `parent` field in the request.
permalink: /233/request-parent-field
redirect_from:
  - /0233/request-parent-field
---

# Batch Create methods: Parent field

This rule enforces that all `BatchCreate` methods have a `string parent` field
in the request message, as mandated in [AIP-233][].

## Details

This rule looks at any message matching `BatchCreate*Request` and complains if
either the `parent` field is missing, or if it has any type other than
`string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message BatchCreateBooksRequest {
  string publisher = 1;  // Field name should be `parent`.
  repeated CreateBookRequest requests = 2;
}
```

```proto
// Incorrect.
message BatchCreateBookRequest {
  bytes parent = 1;  // Field type should be `string`.
  repeated CreateBookRequest requests = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message BatchCreateBookRequest {
  string parent = 1;
  repeated CreateBookRequest requests = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message (if
the `parent` field is missing) or above the field (if it is the wrong type).
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0233::request-parent-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message BatchCreateBooksRequest {
  string publisher = 1;  // Field name should be `parent`.
  repeated CreateBookRequest requests = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-233]: https://aip.dev/233
[aip.dev/not-precedent]: https://aip.dev/not-precedent
