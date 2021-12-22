---
rule:
  aip: 135
  name: [core, '0135', request-name-field]
  summary: Delete RPCs must have a `resource_name` field in the request.
permalink: /135/request-name-field
redirect_from:
  - /0135/request-name-field
---

# Delete methods: Name field

This rule enforces that all `Delete` standard methods have a `string resource_name`
field in the request message, as mandated in [AIP-135][].

## Details

This rule looks at any message matching `Delete*Request` and complains if
either the `resource_name` field is missing, or if it has any type other than `string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message DeleteBookRequest {
  string book = 1;  // Field name should be `resource_name`.
}
```

```proto
// Incorrect.
message DeleteBookRequest {
  bytes resource_name = 1;  // Field type should be `string`.
}
```

**Correct** code for this rule:

```proto
// Correct.
message DeleteBookRequest {
  string resource_name = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message (if
the `resource_name` field is missing) or above the field (if it is the wrong type).
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0135::request-name-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message DeleteBookRequest {
  string book = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-135]: https://aip.dev/135
[aip.dev/not-precedent]: https://aip.dev/not-precedent
