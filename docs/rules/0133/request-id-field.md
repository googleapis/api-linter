---
rule:
  aip: 133
  name: [core, '0133', declarative-friendly-required]
  summary: |
    Declarative-friendly create methods should have a client-specified
    ID field.
permalink: /133/declarative-friendly-required
redirect_from:
  - /0133/declarative-friendly-required
---

# Client-specified IDs

This rule enforces that declarative-friendly create methods include a
client-specified ID field, as mandated in [AIP-133][].

## Details

This rule looks at any `Create` method connected to a resource with a
`google.api.resource` annotation that includes `style: DECLARATIVE_FRIENDLY`,
and complains if it lacks a client-specified ID (e.g. `string book_id`) field.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Assuming that Book is styled declarative-friendly...
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
// Assuming that Book is styled declarative-friendly...
message CreateBookRequest {
  string parent = 1 [(google.api.resource_reference) = {
    child_type: "library.googleapis.com/Book"
  }];

  Book book = 2;

  string book_id = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0133::declarative-friendly-required=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message CreateBookRequest {
  string parent = 1 [(google.api.resource_reference) = {
    child_type: "library.googleapis.com/Book"
  }];

  Book book = 2;
}
```

**Note:** Violations of declarative-friendly rules should be rare, as tools are
likely to expect strong consistency.

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-133]: https://aip.dev/133
[aip.dev/not-precedent]: https://aip.dev/not-precedent
