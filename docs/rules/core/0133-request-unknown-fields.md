---
rule:
  aip: 133
  name: [core, '0133', request-unknown-fields]
  summary: Create RPCs should not have unexpected fields in the request.
---

# Create methods: Unknown fields

This rule enforces that all `Create` standard methods do not have unexpected
fields, as mandated in [AIP-133][].

## Details

This rule looks at any message matching `Create*Request` and complains if it
comes across any fields other than:

- `string parent` ([AIP-133][])
- `{Resource} {resource}` ([AIP-133][])
- `string {resource}_id` ([AIP-133][])
- `string request_id` ([AIP-155][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message CreateBookRequest {
  string parent = 1;
  Book book = 2;
  string book_id = 3;
  string library_id = 4;  // Non-standard field.
}
```

**Correct** code for this rule:

```proto
// Correct.
message CreateBookRequest {
  string parent = 1;
  Book book = 2;
  string book_id = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message CreateBookRequest {
  string parent = 1;
  Book book = 2;
  string book_id = 3;

  // (-- api-linter: core::0133::request-unknown-fields=disabled
  //     aip.dev/not-precedent: We really need this field because reaosns. --)
  string library_id = 4;  // Non-standard field.
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-133]: https://aip.dev/133
[aip-155]: https://aip.dev/155
[aip.dev/not-precedent]: https://aip.dev/not-precedent
