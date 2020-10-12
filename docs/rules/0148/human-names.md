---
rule:
  aip: 148
  name: [core, '0148', human-names]
  summary: Avoid imprecise terms for human names.
permalink: /148/standardized-codes
redirect_from:
  - /0148/standardized-codes
---

# Human names

This rule encourages terms for human names (`given_name` and `family_name`)
that are more accurate across cultures, as mandated in [AIP-148][].

## Details

This rule looks for fields named `first_name` and `last_name`, and complains if
it finds them, suggesting the use of `given_name` and `family_name`
(respectively) instead.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Human {
  string first_name = 1;  // Should be `given_name`.
  string last_name = 2;   // Should be `family_name`
}
```

**Correct** code for this rule:

```proto
// Correct.
message Human {
  string given_name = 1;
  string family_name = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field or its
enclosing message. Remember to also include an [aip.dev/not-precedent][]
comment explaining why.

```proto
// (-- api-linter: core::0148::human-names=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Human {
  string first_name = 1;
  string last_name = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-148]: https://aip.dev/148
[aip.dev/not-precedent]: https://aip.dev/not-precedent
