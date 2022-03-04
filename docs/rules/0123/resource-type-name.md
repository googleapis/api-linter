---
rule:
  aip: 123
  name: [core, '0123', resource-type-name]
  summary: Resource type names must be of the form {Service Name}/{Type}.
permalink: /123/resource-type-name
redirect_from:
  - /0123/resource-type-name
---

# Resource type name

This rule enforces that messages that have a `google.api.resource` annotation,
have a properly formatted `type`, as described in [AIP-123][].

## Details

This rule scans messages with a `google.api.resource` annotation, and validates
the format of the `type` field conforms to `{Service Name}/{Type}`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  option (google.api.resource) = {
    // Should not have more than one separating '/'.
    type: "library.googleapis.com/Genre/Mystery/Book"
    pattern: "publishers/{publisher}/books/{book}"
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
// (-- api-linter: core::0123::resource-type-name=disabled
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
