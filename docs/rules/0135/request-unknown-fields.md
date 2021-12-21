---
rule:
  aip: 135
  name: [core, '0135', request-unknown-fields]
  summary: Delete RPCs should not have unexpected fields in the request.
permalink: /135/request-unknown-fields
redirect_from:
  - /0135/request-unknown-fields
---

# Delete methods: Unknown fields

This rule enforces that all `Delete` standard methods do not have unexpected
fields, as mandated in [AIP-135][].

## Details

This rule looks at any message matching `Delete*Request` and complains if it
comes across any fields other than:

- `string resource_name` ([AIP-135][])
- `bool allow_missing` ([AIP-135][])
- `bool force` ([AIP-135][])
- `string etag` ([AIP-154][])
- `string request_id` ([AIP-155][])
- `bool validate_only` ([AIP-163][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message DeleteBookRequest {
  string name = 1;
  string library_id = 2;  // Non-standard field.
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

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message DeleteBookRequest {
  string resource_name = 1;

  // (-- api-linter: core::0135::request-unknown-fields=disabled
  //     aip.dev/not-precedent: We really need this field because reasons. --)
  string library_id = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-135]: https://aip.dev/135
[aip-154]: https://aip.dev/154
[aip-155]: https://aip.dev/155
[aip-163]: https://aip.dev/163
[aip.dev/not-precedent]: https://aip.dev/not-precedent
