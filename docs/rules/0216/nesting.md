---
rule:
  aip: 216
  name: [core, '0216', nesting]
  summary: Lifecycle state enums should be nested within the resource.
permalink: /216/nesting
redirect_from:
  - /0216/nesting
---

# States

This rule enforces that all lifecycle state enums are nested within their
resource, as recommended in [AIP-216][].

## Details

This rule iterates over enumerations and looks for enums with a name of
`FooState` where a corresponding `Foo` message exists in the same file. If it
finds such a case, it recommends that the enum be nested within the message.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  BookState book_state = 1;
}

// Should be nested under `Book`.
enum BookState {
  BOOK_STATE_UNSPECIFIED = 0;
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  enum State {
    STATE_UNSPECIFIED = 0;
  }
  State state = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the enum.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message Book {
  BookState book_state = 1;
}

// (-- api-linter: core::0216::nesting=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
enum BookState {
  BOOK_STATE_UNSPECIFIED = 0;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-216]: https://aip.dev/216
[aip.dev/not-precedent]: https://aip.dev/not-precedent
