---
rule:
  aip: 214
  name: [core, '0214', ttl-type]
  summary: The `ttl` field should be a `google.protobuf.Duration`.
permalink: /214/ttl-type
redirect_from:
  - /0214/ttl-type
---

# Resource expiry

This rule enforces that `ttl` fields use `google.protobuf.Duration`, as
mandated in [AIP-214][].

## Details

This rule looks at fields named `ttl`, and complains if they are a type other
than `google.protobuf.Duration`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;
  oneof expiration {
    google.protobuf.Timestamp expire_time = 2;
    int32 ttl = 3;  // Should be `google.protobuf.Duration`.
  }
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1;
  oneof expiration {
    google.protobuf.Timestamp expire_time = 2;
    google.protobuf.Duration ttl = 3;
  }
}
```

## Disabling

If you need to violate this rule, use a leading comment above the enum.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message Book {
  string name = 1;
  oneof expiration {
    google.protobuf.Timestamp expire_time = 2;
    // (-- api-linter: core::0214::ttl-type=disabled
    //     aip.dev/not-precedent: We need to do this because reasons. --)
    int32 ttl = 3;  // Should be `google.protobuf.Duration`.
  }
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-214]: https://aip.dev/214
[aip.dev/not-precedent]: https://aip.dev/not-precedent
