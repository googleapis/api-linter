---
rule:
  aip: 152
  name: [core, '0152', request-name-field]
  summary: Run RPCs must have a `name` field in the request.
permalink: /152/request-name-field
redirect_from:
  - /0152/request-name-field
---

# Run requests: Name field

This rule enforces that all `Run` methods have a `string name`
field in the request message, as mandated in [AIP-152][].

## Details

This rule looks at any message matching `Run*JobRequest` and complains if
either the `name` field is missing, or if it has any type other than `string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message RunWriteBookJobRequest {
  string write_book_job = 1;  // Field name should be `name`.
}
```

```proto
// Incorrect.
message RunWriteBookJobRequest {
  bytes name = 1;  // Field type should be `string`.
}
```

**Correct** code for this rule:

```proto
// Correct.
message RunWriteBookJobRequest {
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/WriteBookJob"
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message (if
the `name` field is missing) or above the field (if it is the wrong type).
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0152::request-name-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message RunWriteBookJobRequest {
  string write_book_job = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-152]: https://aip.dev/152
[aip.dev/not-precedent]: https://aip.dev/not-precedent
