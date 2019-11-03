---
rule:
  aip: 143
  name: [core, '0143', string-type]
  summary: Fields representing standardized codes must be strings.
---

# Standardized code strings

This rule attempts to enforce that standard codes for concepts like language,
currency, etc. are strings, as mandated in [AIP-143][].

## Details

This rule looks at any field with a name matching a standardized code, and
complains if it has a type other than `string`.

It currently matches the following field names:

- `currency_code`
- `country_code`
- `language_code`
- `mime_type`
- `time_zone`

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// This enum should not exist.
enum LanguageCode {
  LANGUAGE_CODE_UNSPECIFIED = 0;
  EN_US = 1;
  EN_GB = 2;
}

message Book {
  string name = 1;
  LanguageCode language_code = 2;  // Should be `string`.
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
// (-- api-linter: core::0143::string-type=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
enum LanguageCode {
  LANGUAGE_CODE_UNSPECIFIED = 0;
  EN_US = 1;
  EN_GB = 2;
}

message Book {
  string name = 1;
  LanguageCode language_code = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-143]: https://aip.dev/143
[aip.dev/not-precedent]: https://aip.dev/not-precedent
