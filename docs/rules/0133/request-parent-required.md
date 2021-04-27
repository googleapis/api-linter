---
rule:
  aip: 133
  name: [core, '0133', request-parent-required]
  summary: Create RPCs must have a `parent` field in the request.
permalink: /133/request-parent-required
redirect_from:
  - /0133/request-parent-required
---

# Create methods: Parent field

This rule enforces that all `Create` standard methods have a `string parent`
field in the request message, as mandated in [AIP-133][].

## Details

This rule looks at any message matching `Create*Request` and complains if 
the `parent` field is missing.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message CreateBookRequest {
  // Field name should be `parent`.
  string publisher = 1;
  Book book = 2;
  string book_id = 3;
}
```

**Correct** code for this rule:

```proto
// Correct.
message CreateBookRequest {
  string parent = 1;
  Book book = 2;
  string book_id = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0133::request-parent-required=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message CreateBookRequest {
  string publisher = 1;
  Book book = 2;
  string book_id = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-133]: https://aip.dev/133
[aip.dev/not-precedent]: https://aip.dev/not-precedent
