---
rule:
  aip: 131
  name: [core, '0131', request-has-name-field]
  summary: Get RPCs must have a `name` field in the request.
permalink: /131/request-has-name-field
redirect_from:
  - /0131/request-has-name-field
---

# Get methods: Name field

This rule enforces that all `Get` standard methods have a `string name` field
in the request message, as mandated in [AIP-131][].

## Details

This rule looks at any message matching `Get*Request` and complains if
the `name` field is missing.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message GetBookRequest {
  string book = 1;  // Field name should be `name`.
}
```

**Correct** code for this rule:

```proto
// Correct.
message GetBookRequest {
  string name = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0131::request-has-name-field=disabled
//     aip.dev/not-precedent: This is named "book" for historical reasons. --)
message GetBookRequest {
  string book = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-131]: https://aip.dev/131
[aip.dev/not-precedent]: https://aip.dev/not-precedent
