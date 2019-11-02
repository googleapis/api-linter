---
rule:
  aip: 134
  name: [core, '0134', request-resource-field]
  summary: Update RPCs must have a field for the resource in the request.
---

# Update methods: Resource field

This rule enforces that all `Update` standard methods have a field in the
request message for the resource itself, as mandated in [AIP-134][].

## Details

This rule looks at any message matching `Update*Request` and complains if there
is no field of the resource's type with the expected field name.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// `Book book` is missing.
message UpdateBookRequest {
  google.protobuf.FieldMask update_mask = 2;
}
```

```proto
// Incorrect.
message UpdateBookRequest {
  Book payload = 1;  // Field name should be `book`.
  google.protobuf.FieldMask update_mask = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message UpdateBookRequest {
  Book book = 1;
  google.protobuf.FieldMask update_mask = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message (if
the resource field is missing) or above the field (if it is improperly named).
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0134::request-resource-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message UpdateBookRequest {
  Book payload = 1;
  google.protobuf.FieldMask update_mask = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-134]: https://aip.dev/134
[aip.dev/not-precedent]: https://aip.dev/not-precedent
