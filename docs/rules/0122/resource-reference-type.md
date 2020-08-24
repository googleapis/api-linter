---
rule:
  aip: 122
  name: [core, '0122', resource-reference-type]
  summary: All resource references must be strings.
permalink: /122/resource-reference-type
redirect_from:
  - /0122/resource-reference-type
---

# Resource reference type

This rule enforces that all fields with the `google.api.resource_reference`
annotation are strings, as mandated in [AIP-122][].

## Details

This rule complains if it sees a field with a `google.api.resource_reference`
that has a type other than `string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;

  // Resource references should be strings.
  Author author = 2 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Author"
  }];
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1;

  string author = 2 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Author"
  }];
}
```

```proto
// Correct.
message Book {
  string name = 1;

  // If "author" is not a first-class resource, then it may be a composite
  // field within the book.
  Author author = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message Book {
  string name = 1;

  // (-- api-linter: core::0122::resource-reference-type=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  Author author = 2 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Author"
  }];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-122]: https://aip.dev/122
[aip.dev/not-precedent]: https://aip.dev/not-precedent
