---
rule:
  aip: 123
  name: [core, '0123', resource-name-field]
  summary: Resource messages should have a `string resource_name` field.
permalink: /123/resource-name-field
redirect_from:
  - /0123/resource-name-field
---

# Resource `resource_name` field

This rule enforces that messages that appear to represent resources have a
`string resource_name` field, as described in [AIP-123][].

## Details

This rule scans all messages that have a `google.api.resource` annotation, and
complains if the `resource_name` field is missing or if it is any type other than
singular `string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect: missing `string resource_name` field.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };
}
```

```proto
// Incorrect.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  // Should be `string`, not `bytes`.
  bytes resource_name = 1;
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

  string resource_name = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message, or
above the field if it is the wrong type.

```proto
// (-- api-linter: core::0123::resource-name-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-123]: http://aip.dev/123
[aip.dev/not-precedent]: https://aip.dev/not-precedent
