---
rule:
  aip: 140
  name: [core, '0140', reserved-words]
  summary: Field names must not be reserved words.
permalink: /140/reserved-words
redirect_from:
  - /0140/reserved-words
---

# Field names: Reserved words

This rule enforces that field names are not reserved words, as mandated in
[AIP-140][].

## Details

This rule looks at each field and complains if it the name is a reserved word
in a common programming lanaguge.

Currently, the linter checks all the reserved words in Java, JavaScript, and
Python 3. The exhaustive list of reserved words is found in [the code][].

**Note:** Reserved words in Golang are permitted because Golang's variable
casing rules avoids a conflict.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;
  bool public = 2;  // Reserved word in Java, JavaScript
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1;
  bool is_public = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message Book {
  string name = 1;
  // (-- api-linter: core::0140::reserved-words=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  bool public = 2;  // Reserved word in Java, JavaScript
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

<!-- prettier-ignore-start -->
[aip-140]: https://aip.dev/140
[aip.dev/not-precedent]: https://aip.dev/not-precedent
[the code]: https://github.com/commure/api-linter/blob/main/rules/aip0140/reserved_words.go
<!-- prettier-ignore-end -->
