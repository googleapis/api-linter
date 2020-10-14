---
rule:
  aip: 164
  name: [core, '0164', request-name-field]
  summary: Undelete RPCs must have a `name` field in the request.
permalink: /164/request-name-field
redirect_from:
  - /0164/request-name-field
---

# Undelete methods: Name field

This rule enforces that all `Undelete` methods have a `string name`
field in the request message, as mandated in [AIP-164][].

## Details

This rule looks at any message matching `Undelete*Request` and complains if
either the `name` field is missing, or if it has any type other than `string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message UndeleteBookRequest {
  string book = 1;  // Field name should be `name`.
}
```

```proto
// Incorrect.
message UndeleteBookRequest {
  bytes name = 1;  // Field type should be `string`.
}
```

**Correct** code for this rule:

```proto
// Correct.
message UndeleteBookRequest {
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message (if
the `name` field is missing) or above the field (if it is the wrong type).
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0164::request-name-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message UndeleteBookRequest {
  string book = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-164]: https://aip.dev/164
[aip.dev/not-precedent]: https://aip.dev/not-precedent
