---
rule:
  aip: 163
  name: [core, '0163', synonyms]
  summary: Change validation fields should be named `validate_only`.
permalink: /163/synonyms
redirect_from:
  - /0163/synonyms
---

# Synonyms

This rule enforces that the `validate_only` field is named `validate_only`, and
not a common synonym, as mandated in [AIP-163][].

## Details

This rule complains if it encounters a known synonym to `validate_only`.
Currently, the only recognized synonym `dry_run`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message CreateBookRequest {
  string parent = 1;
  Book book = 2;
  bool dry_run = 3;  // Should be `validate_only`.
}
```

**Correct** code for this rule:

```proto
// Correct.
message CreateBookRequest {
  string parent = 1;
  Book book = 2;
  bool validate_only = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message CreateBookRequest {
  string parent = 1;
  Book book = 2;
  // (-- api-linter: core::0163::synonyms=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  bool dry_run = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-163]: https://aip.dev/163
[aip.dev/not-precedent]: https://aip.dev/not-precedent
