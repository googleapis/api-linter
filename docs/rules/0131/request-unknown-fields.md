---
rule:
  aip: 131
  name: [core, '0131', request-unknown-fields]
  summary: Get RPCs should not have unexpected fields in the request.
permalink: /131/request-unknown-fields
redirect_from:
  - /0131/request-unknown-fields
---

# Get methods: Unknown fields

This rule enforces that all `Get` standard methods do not have unexpected
fields, as mandated in [AIP-131][].

## Details

This rule looks at any message matching `Get*Request` and complains if it comes
across any fields other than:

- `string resource_name` ([AIP-131][])
- `google.protobuf.FieldMask read_mask` ([AIP-157][])
- `View view` ([AIP-157][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message GetBookRequest {
  string resource_name = 1;
  string library_id = 2;  // Non-standard field.
}
```

**Correct** code for this rule:

```proto
// Correct.
message GetBookRequest {
  string resource_name = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message GetBookRequest {
  string resource_name = 1;

  // (-- api-linter: core::0131::request-unknown-fields=disabled
  //     aip.dev/not-precedent: We really need this field because reasons. --)
  string library_id = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-131]: https://aip.dev/131
[aip-157]: https://aip.dev/157
[aip.dev/not-precedent]: https://aip.dev/not-precedent
