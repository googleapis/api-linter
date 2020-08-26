---
rule:
  aip: 140
  name: [core, '0140', numbers]
  summary: Field names should not have words beginning with numbers.
permalink: /140/numbers
redirect_from:
  - /0140/numbers
---

# Field names: Numbers

This rule enforces that field names do not begin any word in the field with a
number, as mandated in [AIP-140][].

## Details

This rule checks every field in the proto and complains if any individual word
begins with a number. It treats the underscore (`_`) character as the only word
separator for this purpose.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;
  int32 review_90th_percentile_stars = 2;
}
```

**Correct** code for this rule:

The correct code here is likely to vary based on the situation. This may be
fixed by spelling out the number:

```proto
// Correct.
message Book {
  string name = 1;
  int32 review_ninetieth_percentile_stars = 2;
}
```

Many cases we see involving numbers like this may be better designed with a
map:

```proto
// Correct.
message Book {
  string name = 1;
  map<int32, int32> review_stars_per_percentile = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0140::numbers=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  string name = 1;
  int32 review_90th_percentile_stars = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-140]: https://aip.dev/140
[aip.dev/not-precedent]: https://aip.dev/not-precedent
