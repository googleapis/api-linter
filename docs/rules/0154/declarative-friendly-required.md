---
rule:
  aip: 154
  name: [core, '0154', declarative-friendly-required]
  summary: Declarative-friendly resources must have an etag field.
permalink: /154/declarative-friendly-required
redirect_from:
  - /0154/declarative-friendly-required
---

# Required etags

This rule enforces that declarative-friendly resources have etags, as mandated
in [AIP-154][].

## Details

This rule looks at any resource with a `google.api.resource` annotation that
includes `style: DECLARATIVE_FRIENDLY`, and complains if it lacks a
`string etag` field.

Additionally, it looks at certain corresponding request messages (e.g.
`DeleteBookRequest`) that _do not_ include the resource, and make the same
check.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
    style: DECLARATIVE_FRIENDLY
  };

  string name = 1;
  // A string etag field should exist.
}
```

```proto
// Incorrect.
message DeleteBookRequest {
  string name = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Book"
  }];
  // A string etag field should exist.
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
    style: DECLARATIVE_FRIENDLY
  }

  string name = 1;
  string etag = 2;
}
```

```proto
// Correct.
message DeleteBookRequest {
  string name = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Book"
  }];
  string etag = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0154::declarative-friendly-required=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
    style: DECLARATIVE_FRIENDLY
  }

  string name = 1;
}
```

**Note:** Violations of declarative-friendly rules should be rare, as tools are
likely to expect strong consistency.

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-154]: https://aip.dev/154
[aip.dev/not-precedent]: https://aip.dev/not-precedent
