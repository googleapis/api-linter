---
rule:
  aip: 145
  name: [core, '0145', start-field]
  summary: Ensures properly named range start field.
permalink: /145/start-field
redirect_from:
  - /0145/start-field
---

# Ranges: Start field

This rule ensures that the field representing the start of a range is properly
named, as described in [AIP-145][].

## Details

This rule complains if it sees a field named `start` without a noun following it
in snake case format.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Chapter {
  // Should have a trailing noun in snake case format.
  int32 start = 1;
}
```

**Correct** code for this rule:

```proto
// Correct.
message Chapter {
  int32 start_page = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0145::start-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Chapter {
  int32 start = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-145]: https://aip.dev/145
[aip.dev/not-precedent]: https://aip.dev/not-precedent
