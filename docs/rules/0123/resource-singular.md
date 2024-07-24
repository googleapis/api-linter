---
rule:
  aip: 123
  name: [core, '0123', resource-singular]
  summary: Resource singular is required and must be lowerCamelCase of type
permalink: /123/resource-singular
redirect_from:
  - /0123/resource-singular
---

# Resource `singular`

This rule enforces that messages that have a `google.api.resource` annotation,
have a properly formatted `singular`, as described in [AIP-123][].

## Details

This rule scans messages with a `google.api.resource` annotation, and validates
the format of the `singular` field is the lower camel case of type.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/BookShelf"
    pattern: "publishers/{publisher}/bookShelves/{book_shelf}"
    // does not match type.
    singular: "shelf",
  };

  string name = 1;
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/BookShelf"
    pattern: "publishers/{publisher}/bookShelves/{book_shelf}"
    singular: "bookShelf",
  };

  string name = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.

```proto
// (-- api-linter: core::0123::resource-singular=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
    singular: "shelf",
  };

  string name = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-123]: http://aip.dev/123
[aip.dev/not-precedent]: https://aip.dev/not-precedent
