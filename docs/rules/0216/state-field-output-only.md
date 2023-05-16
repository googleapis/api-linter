---
rule:
  aip: 216
  name: [core, '0216', state-field-output-only]
  summary: Lifecycle State enum fields should be marked OUTPUT_ONLY
permalink: /216/state-field-output-only
redirect_from:
  - /0216/state-field-output-only
---

# States

This rule requires that all lifecycle state fields whose names end with `State`
are marked as `OUTPUT_ONLY`, as mandated in [AIP-216][].

## Details

This rule iterates over message fields that have an `enum` type, and the type
name ends with `State`. Each field should have the annotation
`[(google.api.field_behavior) = OUTPUT_ONLY]`.

Note that the field name is ignored for the purposes of this rule.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
enum State {  // Should be `State`.
  STATUS_UNSPECIFIED = 0;
}

State state = 1; // Should be marked OUTPUT_ONLY

```

```proto
// Incorrect.
enum BookState {
  BOOK_STATUS_UNSPECIFIED = 0;
  HARDCOVER = 1;
}

BookState state = 1; // Should be marked OUTPUT_ONLY
```

**Correct** code for this rule:

```proto
// Correct.
enum State {
  STATE_UNSPECIFIED = 0;
}

State state = 1 [(google.api.field_behavior) = OUTPUT_ONLY];
```

## Disabling

If you need to violate this rule, use a leading comment above the field
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
enum BookState {
  UNSPECIFIED = 0;
}

// (-- api-linter: core::0216::state-field-output-only=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
BookState state = 1;
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-216]: https://aip.dev/216
[aip.dev/not-precedent]: https://aip.dev/not-precedent
