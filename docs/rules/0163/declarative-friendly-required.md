---
rule:
  aip: 163
  name: [core, '0163', declarative-friendly-required]
  summary: Declarative-friendly mutations should have a validate_only field.
permalink: /163/declarative-friendly-required
redirect_from:
  - /0163/declarative-friendly-required
---

# Required change validation

This rule enforces that declarative-friendly mutations have a `validate_only`
field, as mandated in [AIP-163][].

## Details

This rule looks at any mutation (`POST`, `PATCH`, `DELETE`) connected to a
resource with a `google.api.resource` annotation that includes
`style: DECLARATIVE_FRIENDLY`, and complains if it lacks a `bool validate_only`
field.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Assuming that Book is styled declarative-friendly...
message DeleteBookRequest {
  string name = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Book"
  }];
  // A bool validate_only field should exist.
}
```

**Correct** code for this rule:

```proto
// Correct.
// Assuming that Book is styled declarative-friendly...
message DeleteBookRequest {
  string name = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Book"
  }];
  bool validate_only = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0163::declarative-friendly-required=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message DeleteBookRequest {
  string name = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Book"
  }];
}
```

**Note:** Violations of declarative-friendly rules should be rare, as tools are
likely to expect strong consistency.

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-163]: https://aip.dev/163
[aip.dev/not-precedent]: https://aip.dev/not-precedent
