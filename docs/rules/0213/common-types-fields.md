---
rule:
  aip: 213
  name: [core, '0213', common-types-fields]
  summary: Message fields with certain names should use a common type.
permalink: /213/common-types-fields
redirect_from:
  - /0213/common-types-fields
---

# Common types fields

This rule enforces that message fields with specific names use a common type, as
recommended in [AIP-213][].

## Details

This rule looks at the fields in a message, and complains if it finds a field
with a specific name that doesn't have the corresponding common type.

Some example pairings of common types and field names that are checked are:

* `google.protobuf.Duration`: "duration"
* `google.type.Color`: "color", "colour"
* `google.type.PhoneNumber`: "mobile_number", "phone", "phone_number"

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;
  
  // Should use `google.type.Color`.
  string color = 2;
}
```

```proto
// Incorrect.
message Book {
  string name = 1;
  
  // Should use `google.type.PhoneNumber`.
  string phone_number = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1;
  
  google.type.Color color = 2;
}
```

```proto
// Correct.
message Book {
  string name = 1;
  
  google.type.PhoneNumber phone_number = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the enum.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message Book {
  string name = 1;
  // (-- api-linter: core::0213::common-types-fields=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  google.protobuf.Timestamp expire_time = 2
    [(google.api.field_behavior) = OUTPUT_ONLY];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-213]: https://aip.dev/213
[aip.dev/not-precedent]: https://aip.dev/not-precedent
