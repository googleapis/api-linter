---
rule:
  aip: 202
  name: [core, '0202', string-only-format]
  summary: Certain field formats can only be applied to a string field.
permalink: /202/string-only-format
redirect_from:
  - /0202/string-only-format
---

# String only format

This rule enforces that the following format specifiers are used only on fields
of type `string`, as mandated by [AIP-202][]:

- `UUID4`
- `IPV4`
- `IPV6`
- `IPV4_OR_IPV6`

## Details

This rule looks at every non-string field with a
`(google.api.field_info).format` and complains if the format specifier is one
meant only for use on `string` fields.

## Examples

**Incorrect** code for this rule:

```proto
message Book {
  string name = 1;

  // Incorrect. Non-string must not be assigned format UUID4.
  int32 edition = 2 [(google.api.field_info).format = UUID4];
}
```

**Correct** code for this rule:

```proto
message Book {
  string name = 1;

  int32 edition = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message Book {
  string name = 1;

  // (-- api-linter: core::0202::string-only-format=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  int32 edition = 2 [(google.api.field_info).format = UUID4];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-202]: https://aip.dev/202
[aip.dev/not-precedent]: https://aip.dev/not-precedent