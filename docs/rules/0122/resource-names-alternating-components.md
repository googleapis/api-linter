---
rule:
  aip: 122
  name: [core, '0122', resource-names-alternating-components]
  summary: Resource name components should usually alternate between collection identifiers and resource IDs.
permalink: /122/resource-names-alternating-components
redirect_from:
  - /0122/resource-names-alternating-components
---

# Alternating resource name components

This rule enforces that resource name components should alternate between collection identifiers
(example: `publishers`, `books`, `users`) and resource IDs (example: `123`, `les-miserables`, `vhugo1802`),
as mandated in [AIP-122][].

## Details

This rule scans all `google.api.resource` annotations and checks that the path template segments alternate
between collection identifiers and resource IDs.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    // Should be "books/{book}"
    pattern: "books/book"
  };
  string name = 1;
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };
  string name = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.

```proto
// (-- api-linter: core::0122::resource-names-alternating-components=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/book"
  };
  string name = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-122]: http://aip.dev/122
[aip.dev/not-precedent]: https://aip.dev/not-precedent
