---
rule:
  aip: 122
  name: [core, '0122', resource-collection-identifiers]
  summary: Resource patterns must use lowerCamelCase for collection identifiers.
permalink: /122/resource-collection-identifiers
redirect_from:
  - /0122/resource-collection-identifiers
---

# Resource pattern collection identifiers

This rule enforces that messages that have a `google.api.resource` annotation
have properly formatted collection identifiers in each `pattern`, as described
in [AIP-122][].

## Details

This rule scans messages with a `google.api.resource` annotation, and validates
the format of `pattern` collection identifiers.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    // Collection identifiers must be lowerCamelCase.
    pattern: "Publishers/{publisher}/Books/{book}"
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
    pattern: "publishers/{publisher}/books/{book}"
  };
  string name = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.

```proto
// (-- api-linter: core::0122::resource-collection-identifiers=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "Publishers/{publisher}/Books/{book}"
  };
  string name = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-122]: http://aip.dev/122
[aip.dev/not-precedent]: https://aip.dev/not-precedent
