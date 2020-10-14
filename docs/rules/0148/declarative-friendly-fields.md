---
rule:
  aip: 148
  name: [core, '0148', declarative-friendly-fields]
  summary: Declarative-friendly resources must include some standard fields.
permalink: /148/declarative-friendly-fields
redirect_from:
  - /0148/declarative-friendly-fields
---

# Declarative-friendly fields

This rule requires certain standard fields on declarative-friendly resources,
as mandated in [AIP-148][].

## Details

This rule looks at any resource with a `google.api.resource` annotation that
includes `style: DECLARATIVE_FRIENDLY`, and complains if it does not include
all of the following fields:

- `string name`
- `string uid`
- `string display_name`
- `google.protobuf.Timestamp create_time`
- `google.protobuf.Timestamp update_time`
- `google.protobuf.Timestamp delete_time`

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}
    style: DECLARATIVE_FRIENDLY
  };

  string name = 1;
  // string uid should be included!
  string display_name = 2;
  google.protobuf.Timestamp create_time = 3;
  google.protobuf.Timestamp update_time = 4;
  // google.protobuf.TImestamp delete_time should be included!
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}
    style: DECLARATIVE_FRIENDLY
  };

  string name = 1;
  string uid = 2;
  string display_name = 3;
  google.protobuf.Timestamp create_time = 4;
  google.protobuf.Timestamp update_time = 5;
  google.protobuf.TImestamp delete_time = 6;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0148::declarative-friendly-fields=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}
    style: DECLARATIVE_FRIENDLY
  };

  string name = 1;
  // string uid should be included!
  string display_name = 2;
  google.protobuf.Timestamp create_time = 3;
  google.protobuf.Timestamp update_time = 4;
  // google.protobuf.TImestamp delete_time should be included!
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-148]: https://aip.dev/148
[aip.dev/not-precedent]: https://aip.dev/not-precedent
