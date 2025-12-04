---
rule:
  aip: 133
  name: [core, '0133', request-id-field]
  summary: create methods should have a client-specified ID field.
permalink: /133/request-id-field
redirect_from:
  - /0133/request-id-field
---

# Client-specified IDs

This rule enforces that declarative-friendly create methods include a
client-specified ID field, as mandated in [AIP-133][].

## Details

This rule looks at any `Create` method connected to a declarative-friendly
resource (one with `style: DECLARATIVE_FRIENDLY` in its `google.api.resource`
annotation) and complains if it lacks a client-specified ID (e.g. `string
book_id`) field.

**Note:** This rule only applies to declarative-friendly resources. Resources
without the `DECLARATIVE_FRIENDLY` style are not checked by this rule.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Book is a declarative-friendly resource.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
    style: DECLARATIVE_FRIENDLY
  };
}

message CreateBookRequest {
  string parent = 1 [(google.api.resource_reference) = {
    child_type: "library.googleapis.com/Book"
  }];

  Book book = 2;

  // A `string book_id` field should exist.
}
```

**Correct** code for this rule:

```proto
// Correct.
// Book is a declarative-friendly resource.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
    style: DECLARATIVE_FRIENDLY
  };
}

message CreateBookRequest {
  string parent = 1 [(google.api.resource_reference) = {
    child_type: "library.googleapis.com/Book"
  }];

  string book_id = 2;

  Book book = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0133::request-id-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message CreateBookRequest {
  string parent = 1 [(google.api.resource_reference) = {
    child_type: "library.googleapis.com/Book"
  }];

  Book book = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-133]: https://aip.dev/133
[aip.dev/not-precedent]: https://aip.dev/not-precedent