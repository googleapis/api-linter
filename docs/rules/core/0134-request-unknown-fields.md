---
rule:
  aip: 134
  name: [core, '0134', request-unknown-fields]
  summary: Update RPCs should not have unexpected fields in the request.
---

# Update methods: Unknown fields

This rule enforces that all `Update` standard methods do not have unexpected
fields, as mandated in [AIP-134][].

## Details

This rule looks at any message matching `Update*Request` and complains if it
comes across any fields other than:

- `{Resource} {resource}` ([AIP-134][])
- `google.protobuf.FieldMask update_mask` ([AIP-134][])
- `string request_id` ([AIP-155][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message UpdateBookRequest {
  Book book = 1;
  google.protobuf.FieldMask update_mask = 2;
  string library_id = 3;  // Non-standard field.
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

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message UpdateBookRequest {
  Book book = 1;
  google.protobuf.FieldMask update_mask = 2;
  // (-- api-linter: core::0134::request-unknown-fields=disabled
  //     aip.dev/not-precedent: We really need this field because reaosns. --)
  string library_id = 3;  // Non-standard field.
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-134]: https://aip.dev/134
[aip-155]: https://aip.dev/155
[aip.dev/not-precedent]: https://aip.dev/not-precedent
