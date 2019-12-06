---
rule:
  aip: 123
  name: [core, '0123', resource-variables]
  summary: Resource patterns should use consistent variable naming.
permalink: /123/resource-variables
redirect_from:
  - /0123/resource-variables
---

# Resource pattern variables

This rule enforces that resource patterns use consistent variable naming
conventions, as described in [AIP-123][].

## Details

This rule scans all messages with `google.api.resource` annotations, and
complains if variables in a `pattern` use camel case, or end in `_id`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    // Should be: publishers/{publisher}/books/{book}
    pattern: "publishers/{publisher_id}/books/{book_id}"
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
    pattern: "publishers/{publisher}/electronicBooks/{electronicBook}"
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
// (-- api-linter: core::0123::resource-variables=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher_id}/books/{book_id}"
  }

  string name = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-123]: http://aip.dev/123
[aip.dev/not-precedent]: https://aip.dev/not-precedent
