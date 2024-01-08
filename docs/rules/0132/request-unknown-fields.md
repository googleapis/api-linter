---
rule:
  aip: 132
  name: [core, '0132', request-unknown-fields]
  summary: List RPCs should not have unexpected fields in the request.
permalink: /132/request-unknown-fields
redirect_from:
  - /0132/request-unknown-fields
---

# List methods: Unknown fields (Request)

This rule enforces that all `List` standard methods do not have unexpected
fields, as mandated in [AIP-132][].

## Details

This rule looks at any message matching `List*Request` and complains if it
comes across any fields other than:

- `string parent` ([AIP-132][])
- `int32 page_size` ([AIP-158][])
- `string page_token` ([AIP-158][])
- `int32 skip` ([AIP-158][])
- `string filter` ([AIP-132][])
- `string order_by` ([AIP-132][])
- `bool show_deleted` ([AIP-132][])
- `string request_id` ([AIP-155][])
- `google.protobuf.FieldMask read_mask` ([AIP-157][])
- `View view` ([AIP-157][])

It only checks field names; it does not validate type correctness. This is
handled by other rules, such as
[request field types](./0132-request-field-types.md).

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message ListBooksRequest {
  string parent = 1;
  int32 page_size = 2;
  string page_token = 3;
  string library_id = 4;  // Non-standard field.
}
```

**Correct** code for this rule:

```proto
// Correct.
message ListBooksRequest {
  string parent = 1;
  int32 page_size = 2;
  string page_token = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message ListBooksRequest {
  string parent = 1;
  int32 page_size = 2;
  string page_token = 3;

  // (-- api-linter: core::0132::request-unknown-fields=disabled
  //     aip.dev/not-precedent: We really need this field because reasons. --)
  string library_id = 4;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-132]: https://aip.dev/132
[aip-135]: https://aip.dev/135
[aip-155]: https://aip.dev/155
[aip-157]: https://aip.dev/157
[aip-158]: https://aip.dev/158
[aip.dev/not-precedent]: https://aip.dev/not-precedent
