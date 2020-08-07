---
rule:
  aip: 124
  name: [core, '0124', valid-reference]
  summary: Resource patterns should use consistent variable naming.
permalink: /124/valid-reference
redirect_from:
  - /0124/valid-reference
---

# Valid resource references

This rule enforces that resource reference annotations refer to valid and
reachable resource types, as described in [AIP-124][].

## Details

This rule scans all fields with `google.api.resource_reference` annotations,
and complains if the `type` on them refers to a resource with no corresponding
`google.api.resource` or `google.api.resource_definition`.

The rule scans the file where the field is found and all files imported by that
file (recursively) as long as they are in the same package.

Certain common resource types are exempt from this rule.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  // Needs a resource annotation; without one, resource references are invalid.
  string name = 1;

  // ...
}

message GetBookRequest {
  string name = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Book"  // Lint warning; reference not found.
  }]
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  string name = 1;

  // ...
}

message GetBookRequest {
  string name = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Book"
  }]
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message Book {
  string name = 1;
}

message GetBookRequest {
  // (-- api-linter: core::0124::valid-reference=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string name = 1 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Book"
  }]
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-124]: http://aip.dev/124
[aip.dev/not-precedent]: https://aip.dev/not-precedent

```

```
