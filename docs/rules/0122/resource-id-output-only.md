---
rule:
  aip: 122
  name: [core, '0122', resource-id-output-only]
  summary: Resource ID fields must be classified as `OUTPUT_ONLY`.
permalink: /122/resource-id-output-only
redirect_from:
  - /0122/resource-id-output-only
---

# Output Only Resource ID fields

This rule enforces that resource ID fields are classified as `OUTPUT_ONLY`, as
mandated in [AIP-122][].

## Details

This rule scans all resource fields and complains if it sees an ID field, named
as `uid` or with the `_id` suffix, that is not classified as `OUTPUT_ONLY`.

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
  // Should have `(google.api.field_behavior) = OUTPUT_ONLY`.
  string book_id = 2;
  // Should have `(google.api.field_behavior) = OUTPUT_ONLY`.
  string uid = 3;
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
  string book_id = 2 [(google.api.field_behavior) = OUTPUT_ONLY];
  string uid = 3 [(google.api.field_behavior) = OUTPUT_ONLY];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.

```proto
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };
  string name = 1;
  // (-- api-linter: core::0122::resource-id-output-only=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string book_id = 2;
  // (-- api-linter: core::0122::resource-id-output-only=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string uid = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-122]: http://aip.dev/122
[aip.dev/not-precedent]: https://aip.dev/not-precedent
