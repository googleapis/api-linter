---
rule:
  aip: 128
  name: [core, '0128', resource-reconciling-behavior]
  summary: Declarative-friendly resources should annotate the `reconciling` field as `OUTPUT_ONLY`.
permalink: /128/resource-reconciling-behavior
redirect_from:
  - /0128/resource-reconciling-behavior
---

# Declarative-friendly resources: Reconciling field behavior

This rule enforces that all declarative-friendly resources have
`google.api.field_behavior` set to `OUTPUT_ONLY` on their `bool
reconciling` field, as mandated in [AIP-128][].

## Details

This rule looks at any message with a `google.api.resource` annotation that
includes `style: DECLARATIVE_FRIENDLY`, and complains if the `reconciling` field
does not have a `google.api.field_behavior` annotation with a value of
`OUTPUT_ONLY`.

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

  // The `google.api.field_behavior` annotation should be `OUTPUT_ONLY`.
  bool reconciling = 2;
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

  bool reconciling = 2 [(google.api.field_behavior) = OUTPUT_ONLY];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
    style: DECLARATIVE_FRIENDLY
  };

  string name = 1;

  // (-- api-linter: core::0128::resource-reconciling-behavior=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  bool reconciling = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-128]: https://aip.dev/128
[aip.dev/not-precedent]: https://aip.dev/not-precedent
