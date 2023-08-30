---
rule:
  aip: 203
  name: [core, '0203', resource-identifier-only]
  summary: Only the resource's name field can have field behavior IDENTIFIER.
permalink: /203/resource-identifier-only
redirect_from:
  - /0203/resource-identifier-only
---

# Resource name: IDENTIFIER

This rule enforces that only the field representing the resource's name is
annotated with with `(google.api.field_behavior) = IDENTIFIER`, as mandated
by [AIP-203][].

## Details

This rule looks at every field with `(google.api.field_behavior) = IDENTIFER`
and complains if that field is not the resource's name field.

## Examples

**Incorrect** code for this rule:

```proto
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };
  string name = 1 [(google.api.field_behavior) = IDENTIFER];

  // Incorrect. Must not be assigned IDENTIFIER.
  string description = 2 [(google.api.field_behavior) = IDENTIFER];
}
```

**Correct** code for this rule:

```proto
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };
  string name = 1 [(google.api.field_behavior) = IDENTIFIER];

  // Correct.
  string description = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };
  string name = 1 [(google.api.field_behavior) = IDENTIFER];

  // (-- api-linter: core::0203::resource-identifier-only=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string description = 2 [(google.api.field_behavior) = IDENTIFER];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-203]: https://aip.dev/203
[aip.dev/not-precedent]: https://aip.dev/not-precedent