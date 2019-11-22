---
rule:
  aip: 132
  name: [core, '0132', request-unknown-fields]
  summary: List RPCs should not have unexpected fields in the request.
permalink: /132/request-unknown-fields
redirect_from:
  - /0132/request-unknown-fields
---

# List methods: Unknown fields (Response)

This rule enforces that all `List` standard methods do not have unexpected
fields, as mandated in [AIP-132][].

## Details

This rule looks at any message matching `List*Request` and complains if it
comes across any fields other than:

- The resource.
- `int32/int64 total_size` ([AIP-132][])
- `string next_page_token` ([AIP-158][])
- `repeated string unavailable` ([AIP-217][])

It only checks field names; it does not validate type correctness or
repeated-ness.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message ListBooksResponse {
  repeated Book books = 1;
  string next_page_token = 2;
  string publisher_id = 3;  // Unrecognized field.
}
```

**Correct** code for this rule:

```proto
// Correct.
message ListBooksResponse {
  repeated Book books = 1;
  string next_page_token = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message ListBooksResponse {
  repeated Book books = 1;
  string next_page_token = 2;
  // (-- api-linter: core::0132::response-unknown-fields=disabled
  //     aip.dev/not-precedent: We really need this field because reasons. --)
  string publisher_id = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-132]: https://aip.dev/132
[aip-135]: https://aip.dev/135
[aip-157]: https://aip.dev/157
[aip-158]: https://aip.dev/158
[aip.dev/not-precedent]: https://aip.dev/not-precedent
