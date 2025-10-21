---
rule:
  aip: 123
  name: [core, '0123', resource-pattern-plural]
  summary: Resource patterns must use the plural as the collection segment
permalink: /123/resource-pattern-plural
redirect_from:
  - /0123/resource-pattern-plural
---

# Resource `pattern` use of `plural`

This rule enforces that messages that have a `google.api.resource` annotation
use the `plural` form as the collection segment, as described in [AIP-123][].

## Details

This rule scans messages with a `google.api.resource` annotation, and validates
the `plural` form of the resource type name is used as the collection segment
in every pattern.

**Note:** Special consideration is given to type names of child collections that
stutter next to their parent collection, as described in
[AIP-122 Nested Collections][nested]. See AIP-122 for more details.

**Important:** Do not accept the suggestion if it would produce a backwards
incompatible change.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/BookShelf"
    // collection segment doesn't match the plural.
    pattern: "publishers/{publisher}/shelves/{book_shelf}"
    singular: "bookShelf"
    plural: "bookShelves"
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
    singular: "bookShelf"
    plural: "bookShelves"
  };

  string name = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.

```proto
// (-- api-linter: core::0123::resource-pattern-plural=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/BookShelf"
    pattern: "publishers/{publisher}/shelves/{book_shelf}"
    singular: "bookShelf"
    plural: "bookShelves"
  };

  string name = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-123]: http://aip.dev/123
[aip.dev/not-precedent]: https://aip.dev/not-precedent
[nested]: https://aip.dev/122#nested-collections
