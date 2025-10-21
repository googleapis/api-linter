---
rule:
  aip: 148
  name: [core, '0148', use-uid]
  summary: Use uid instead of id in resource messages.
permalink: /148/use-uid
redirect_from:
  - /0148/use-uid
---

# Use `uid` as the resource ID field

This rule encourages the use of `uid` instead of `id` for resource messages, as
mandated in [AIP-148][].

## Details

This rule looks on resource messages for a field named `id`, complains if it
finds it, and suggests the use of `uid` instead.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };

  string name = 1;
  string id = 2; // Should be `uid`
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };

  string name = 1;
  string uid = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field or its
enclosing message. Remember to also include an [aip.dev/not-precedent][]
comment explaining why.

```proto
// (-- api-linter: core::0148::use-uid=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };

  string name = 1;
  string id = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-148]: https://aip.dev/148
[aip.dev/not-precedent]: https://aip.dev/not-precedent
