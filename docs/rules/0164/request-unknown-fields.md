---
rule:
  aip: 164
  name: [core, '0164', request-unknown-fields]
  summary: Undelete RPCs should not have unexpected fields in the request.
permalink: /164/request-unknown-fields
redirect_from:
  - /0164/request-unknown-fields
---

# Undelete methods: Unknown fields

This rule enforces that all `Undelete` requests do not have unexpected
fields, as mandated in [AIP-164][].

## Details

This rule looks at any message matching `Undelete*Request` and complains if it
comes across any fields other than:

- `string name` ([AIP-164][])
- `string etag` ([AIP-154][])
- `string request_id` ([AIP-155][])
- `bool validate_only` ([AIP-163][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message UndeleteBookRequest {
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book",
  ];
  string library_id = 2;  // Non-standard field.
}
```

**Correct** code for this rule:

```proto
// Correct.
message UndeleteBookRequest {
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book",
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message UndeleteBookRequest {
  string name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book",
  ];

  // (-- api-linter: core::0164::request-unknown-fields=disabled
  //     aip.dev/not-precedent: We really need this field because reasons. --)
  string library_id = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-154]: https://aip.dev/154
[aip-155]: https://aip.dev/155
[aip-163]: https://aip.dev/163
[aip-164]: https://aip.dev/164
[aip.dev/not-precedent]: https://aip.dev/not-precedent
