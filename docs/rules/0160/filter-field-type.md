---
rule:
  aip: 160
  name: [core, '0160', filter-field-type]
  summary: |
    The filtering field on List and custom method request messages,
    "filter" must be a string.
permalink: /160/filter-field-type
redirect_from:
  - /0160/filter-field-type
---

# Filtering: filter field type

This rule enforces that the field used for filtering, called `filter`,
is of type `string`, as mandated by [AIP-160][].

## Details

This rule looks at the request message for List methods and custom
methods. Then, it verifes that any field called `filter` is of type
`string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {
  bytes filter = 1;
}
```

```proto
// Incorrect.
rpc CheckoutBook(CheckoutBookRequest) returns (CheckoutBookResponse) {
  BookFilters filter = 1;

  message BookFilters {
    string publisher = 1;
    repeated string authors = 2;
  }
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {
  string filter = 1;
}
```

```proto
// Correct.
rpc CheckoutBook(CheckoutBookRequest) returns (CheckoutBookResponse) {
  string filter = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0160::filter-field-type=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {
  bytes filter = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-160]: https://aip.dev/160
[aip.dev/not-precedent]: https://aip.dev/not-precedent

