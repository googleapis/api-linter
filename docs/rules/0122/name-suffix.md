---
rule:
  aip: 122
  name: [core, '0122', name-suffix]
  summary: Fields should not use the suffix `_name`.
permalink: /122/name-suffix
redirect_from:
  - /0122/name-suffix
---

# Name field suffix

This rule enforces that fields do not use the suffix `_name`, as mandated in
[AIP-122][].

## Details

This rule scans all fields complains if it sees the suffix `_name` on a field.

**Note:** The standard field `display_name` is exempt, as are `given_name` and
`family_name` (used to represent the human-readable name of a person),
`full_resource_name`, and `crypto_key_name`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;
  string author_name = 2;  // Should be `author`.
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1;
  string author = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.

```proto
// (-- api-linter: core::0122::name-suffix=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  string name = 1;
  string author_name = 2;  // Should be `author`.
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-122]: http://aip.dev/122
[aip.dev/not-precedent]: https://aip.dev/not-precedent
