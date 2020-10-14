---
rule:
  aip: 128
  name: [core, '0128', resource-reconciling-field]
  summary: Declarative-friendly resources must have a `reconciling` field.
permalink: /128/resource-reconciling-field
redirect_from:
  - /0128/resource-reconciling-field
---

# Declarative-friendly resources: Reconciling field

This rule enforces that all declarative-friendly resources have a `bool
reconciling` field, as mandated in [AIP-128][].

## Details

This rule looks at any message with a `google.api.resource` annotation that
includes `style: DECLARATIVE_FRIENDLY`, and complains if the `reconciling` field
is missing or if it has any type other than `bool`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// The `reconciling` field is missing.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
    style: DECLARATIVE_FRIENDLY
  };

  string name = 1;
}
```

```proto
// Incorrect.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
    style: DECLARATIVE_FRIENDLY
  };

  string name = 1;
  int32 reconciling = 2; // Type should be `bool`.
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
  };

  string name = 1;
  bool reconciling = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message (if
the `reconciling` field is missing) or above the field (if it is the wrong
type). Remember to also include an [aip.dev/not-precedent][] comment explaining
why.

```proto
// (-- api-linter: core::0128::resource-reconciling-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
    style: DECLARATIVE_FRIENDLY
  };

  string name = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-128]: https://aip.dev/128
[aip.dev/not-precedent]: https://aip.dev/not-precedent
