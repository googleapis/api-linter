---
rule:
  aip: 143
  name: [core, '0143', standard-codes]
  summary: Fields representing concepts with standardized codes must use them.
permalink: /143/standardized-codes
redirect_from:
  - /0143/standardized-codes
---

# Standardized codes

This rule attempts to enforce that standard codes for concepts like language,
currency, etc. are consistently used rather than any alternatives, as mandated
in [AIP-143][].

## Details

This rule looks at any field with a name that looks close to a field with a
common standardized code, but that is not exactly that. It complains if it
finds one and suggests the correct field name.

It currently spots the following common substitutes:

- `content_type`
- `country`
- `currency`
- `lang`
- `language`
- `mime`
- `mimetype`
- `tz`
- `timezone`

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;
  string lang = 2;  // Should be `language_code`.
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1;
  string language_code = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0143::standard-codes=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  string name = 1;
  string lang = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-143]: https://aip.dev/143
[aip.dev/not-precedent]: https://aip.dev/not-precedent
