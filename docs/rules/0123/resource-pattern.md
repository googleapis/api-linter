---
rule:
  aip: 123
  name: [core, '0123', resource-pattern]
  summary: Resource annotations should define a pattern.
permalink: /123/resource-pattern
redirect_from:
  - /0123/resource-pattern
---

# Resource patterns

This rule enforces that messages that appear to represent resources have a
`pattern` defined on their `google.api.resource` annotation, as described in
[AIP-123][].

## Details

This rule scans all messages with `google.api.resource` annotations, and
complains if `pattern` is not provided at least once. It also complains if the
segments outside of variable names contain underscores.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    // pattern should be here
  }

  string name = 1;
}
```

```proto
// Incorrect.
message ElectronicBook {
  option (google.api.resource) = {
    type: "library.googleapis.com/ElectronicBook"
    // Should be: publishers/{publisher}/electronicBooks/{electronic_book}
    pattern: "publishers/{publisher}/electronic_books/{electronic_book}"
  }

  string name = 1;
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  string name = 1;
}
```

```proto
// Correct.
message ElectronicBook {
  option (google.api.resource) = {
    type: "library.googleapis.com/ElectronicBook"
    pattern: "publishers/{publisher}/electronicBooks/{electronic_book}"
  };

  string name = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.

```proto
// (-- api-linter: core::0123::resource-pattern=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
  }

  string name = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-123]: http://aip.dev/123
[aip.dev/not-precedent]: https://aip.dev/not-precedent
