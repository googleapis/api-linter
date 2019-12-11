---
rule:
  aip: 132
  name: [core, '0132', request-parent-required]
  summary: List RPCs must have a `parent` field in the request.
permalink: /132/request-parent-required
redirect_from:
  - /0132/request-parent-required
---

# List methods: Parent field

This rule enforces that all `List` standard methods have a `string parent`
field in the request message, as mandated in [AIP-132](http://aip.dev/132).

## Details

This rule looks at any message matching `List*Request` and complains if
the `parent` field is missing.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message ListBooksRequest {
  // Field name should be `parent`.
  string publisher = 1;  
  int32 page_size = 2;
  string page_token = 3;
}
```

**Correct** code for this rule:

```proto
// Correct.
message ListBooksRequest {
  string parent = 1;
  int32 page_size = 2;
  string page_token = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0132::request-parent-required=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message ListBooksRequest {
  string publisher = 1;
  int32 page_size = 2;
  string page_token = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip.dev/not-precedent]: https://aip.dev/not-precedent
