---
rule:
  aip: 131
  name: [core, '0131', request-name-field]
  summary: Get RPCs must have a `name` field in the request.
---

# Get methods: Name field

This rule enforces that all `Get` standard methods have a `string name` field
in the request message, as mandated in [AIP-131](http://aip.dev/131).

## Details

This rule looks at any message matching `Get*Request` and complains if either
the `name` field is missing, or if it has any type other than `string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message GetBookRequest {
  string book = 1;  // Field name should be `name`.
}
```

```proto
// Incorrect.
message GetBookRequest {
  bytes name = 1;  // Field type should be `string`.
}
```

**Correct** code for this rule:

```proto
// Correct.
message GetBookRequest {
  string name = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message (if
the `name` field is missing) or above the field (if it is the wrong type).

```proto
// (-- api-linter: core::0131::request-name-field=disabled --)
message GetBookRequest {
  string book = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.
