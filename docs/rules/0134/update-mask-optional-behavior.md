---
rule:
  aip: 134
  name: [core, '0134', update-mask-optional-behavior]
  summary: Standard Update `update_mask` field must be `OPTIONAL`.
permalink: /134/update-mask-optional-behavior
redirect_from:
  - /0134/update-mask-optional-behavior
---

# Update methods: Mask field expected behavior

This rule enforces that the `update_mask` field of a Standard `Update` request
uses `google.api.field_behavior = OPTIONAL`, as mandated in [AIP-134][].

## Details

This rule looks at any field named `update_mask` that's in a `Update*Request`
and complains if it the field is not annotated with
`google.api.field_behavior = OPTIONAL`.

## Examples

**Incorrect** code for this rule:

```proto
message UpdateBookRequest {
  Book book = 1;

  // Incorrect. Must be `OPTIONAL`.
  google.protobuf.FieldMask update_mask = 2 [  
    (google.api.field_behavior) = REQUIRED
  ];
}
```


**Correct** code for this rule:

```proto
message UpdateBookRequest {
  Book book = 1;

  // Correct.
  google.protobuf.FieldMask update_mask = 2 [
    (google.api.field_behavior) = OPTIONAL
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message UpdateBookRequest {
  Book book = 1;

  // (-- api-linter: core::0134::update-mask-optional-behavior=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  google.protobuf.FieldMask update_mask = 2 [  
    (google.api.field_behavior) = REQUIRED
  ];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-134]: https://aip.dev/134
[aip.dev/not-precedent]: https://aip.dev/not-precedent
