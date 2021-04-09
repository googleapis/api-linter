--------------------------------------------------------------------------------

rule: aip: 192 name: [core, '0192', deprecated-service-comment] summary:
Deprecated services must have a corresponding comment. permalink:
/192/deprecated-service-comment redirect_from:

## - /0192/deprecated-service-comment

# Omnipresent comments

This rule enforces that every service marked with the protobuf
`deprecated`option has `"Deprecated: <reason>"` as the first line in the public
leading comment, as mandated in [AIP-192][].

## Details

This rule looks at each service descriptor in each proto file, and complains if
the protobuf `deprecated` option is set, but the first line of the public
comment does not begin with "Deprecated: ".

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// A library service.
service Library {
  option deprecated = true;
  rpc GetBook(GetBookRequest) returns (Book);
}
```

**Correct** code for this rule:

```proto
// Deprecated: Please visit a your local library instead.
// Correct.
// A library service.
service Library {
  option deprecated = true;
  rpc GetBook(GetBookRequest) returns (Book);
}
```

## Disabling

If you need to violate this rule, use a leading comment above the descriptor.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// A representation of a book.
message Book {
  // (-- api-linter: core::0192::deprecated-service-comment=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string name = 1;
}
```

[aip-192]: https://aip.dev/192
[aip.dev/not-precedent]: https://aip.dev/not-precedent
