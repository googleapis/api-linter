---
rule:
  aip: 123
  name: [core, '0123', resource-plural]
  summary: Resource plural is required
permalink: /123/resource-plural
redirect_from:
  - /0123/resource-plural
---

# Resource type name

This rule enforces that messages that have a `google.api.resource` annotation,
have a properly formatted `plural`, as described in [AIP-123][].

## Details

This rule scans messages with a `google.api.resource` annotation, and
verifies the `plural` field exists.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  // no plural annotation
  option (google.api.resource) = {
    type: "library.googleapis.com/BookShelf"
    pattern: "publishers/{publisher}/bookShelves/{book_shelf}"
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
    plural: "bookShelves",
  };

  string name = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.

```proto
// (-- api-linter: core::0123::resource-plural=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Genre/Mystery/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  string name = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-123]: http://aip.dev/123
[aip.dev/not-precedent]: https://aip.dev/not-precedent
