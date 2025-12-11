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

This rule enforces that create methods include a client-specified ID field, as
mandated in [AIP-133][].

## Details

This rule looks at any `Create` method connected to a resource and complains if
it lacks a client-specified ID (e.g. `string book_id`) field.

**Note:** This rule does not apply if the request contains a `request_id` field.
This accommodates [AIP-133 exception cases][user-specified-ids] for data plane
resources that may use `request_id` for idempotency ([AIP-155][]) without
requiring a client-specified resource ID.

[user-specified-ids]: https://google.aip.dev/133#user-specified-ids
[aip-155]: https://google.aip.dev/155

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
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
