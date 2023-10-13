---
rule:
  aip: 148
  name: [core, '0148', ip-address-format]
  summary: Annotate IP address fields with an IP address format.
permalink: /148/ip-address-format
redirect_from:
  - /0148/ip-address-format
---

# IP Address field format annotation

This rule encourages the use of one of the IP Address format annotations,
`IPV4`, `IPV6`, or `IPV4_OR_IPV6`, on the `ip_address` field or a field ending
with `_ip_address`, as mandated in [AIP-148][].

## Details

This rule looks on for fields named `ip_address` or ending with `_ip_address`
and complains if it does not have the `(google.api.field_info).format`
annotation with one of `IPV4`, `IPV6`, or `IPV4_OR_IPV6`, or has a format other
than than one of those.

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
  string ip_address = 2; // missing (google.api.field_info).format = IPV4
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
  string ip_address = 2 [(google.api.field_info).format = IPV4];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field or its
enclosing message. Remember to also include an [aip.dev/not-precedent][]
comment explaining why.

```proto
// (-- api-linter: core::0148::ip-address-format=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };

  string name = 1;

  string ip_address = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-148]: https://aip.dev/148
[aip.dev/not-precedent]: https://aip.dev/not-precedent