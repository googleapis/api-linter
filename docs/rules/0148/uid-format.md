---
rule:
  aip: 148
  name: [core, '0148', uid-format]
  summary: Annotate uid with UUID4 format.
permalink: /148/uid-format
redirect_from:
  - /0148/uid-format
---

# `uid` format annotation

This rule encourages the use of the `UUID4` format annotation on the `uid`
field, as mandated in [AIP-148][].

## Details

This rule looks on for fields named `uid` and complains if it does not have the
`(google.api.field_info).format = UUID4` annotation or has a
[format](https://github.com/googleapis/googleapis/blob/master/google/api/field_info.proto#L54)
other than `UUID4`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };

  string name = 1 [(google.api.field_behavior) = IDENTIFIER];
  string uid = 2; // missing (google.api.field_info).format = UUID4
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

  string name = 1 [(google.api.field_behavior) = IDENTIFIER];
  string uid = 2 [(google.api.field_info).format = UUID4];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field or its
enclosing message. Remember to also include an [aip.dev/not-precedent][]
comment explaining why.

```proto
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };

  string name = 1 [(google.api.field_behavior) = IDENTIFIER];

  // (-- api-linter: core::0148::uid-format=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string uid = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-148]: https://aip.dev/148
[aip.dev/not-precedent]: https://aip.dev/not-precedent
