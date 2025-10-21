---
rule:
  aip: 122
  name: [core, '0122', no-self-links]
  summary: Resources should not contain self-links.
permalink: /122/no-self-links
redirect_from:
  - /0122/no-self-links
---

# No self links

This rule enforces that resource messages do not contain any fields called
`string self_link`, as mandated in [AIP-122][].

## Details

This rule complains if it sees a resource field of type `string` that is also
named `self_link`.

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

  // Incorrect. Resources should contain self-links.
  string self_link = 2;
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
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// Incorrect.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };
  string name = 1;

  // (-- api-linter: core::0122::no-self-links=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string self_link = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-122]: https://aip.dev/122
[aip.dev/not-precedent]: https://aip.dev/not-precedent
