---
rule:
  aip: 131
  name: [core, '0131', request-name-required]
  summary: Get RPCs must have a `resource_name` field in the request.
permalink: /131/request-name-required
redirect_from:
  - /0131/request-name-required
---

# Get methods: Resource Name field

This rule enforces that all `Get` standard methods have a `string resource_name` field
in the request message, as mandated in [AIP-131][].

## Details

This rule looks at any message matching `Get*Request` and complains if
the `resource_name` field is missing.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message GetBookRequest {
  string book = 1 [  // Field name should be `resource_name`.
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

**Correct** code for this rule:

```proto
// Correct.
message GetBookRequest {
  string resource_name = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0131::request-name-required=disabled
//     aip.dev/not-precedent: This is named "book" for historical reasons. --)
message GetBookRequest {
  string book = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-131]: https://aip.dev/131
[aip.dev/not-precedent]: https://aip.dev/not-precedent
