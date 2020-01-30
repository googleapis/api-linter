---
rule:
  aip: 136
  name: [core, '0136', http-uri-suffix]
  summary: Custom methods should have a correct URI suffix.
permalink: /136/http-uri-suffix
redirect_from:
  - /0136/http-uri-suffix
---

# Custom methods: URI suffix

This rule enforces that custom methods include the custom verb in the REST URI,
as mandated in [AIP-136][].

## Details

This rule looks at any method that is not a standard method, and tries to find
the appropriate suffix at the end of the URI. More specifically:

- If the URI contains a `name` or `parent` variable, then it expects `:verb` at
  the end of the URI.
- Otherwise, it expects `:verbNoun` at the end of the URI.

**Note:** This rule will not run if the [http-name-variable][] or
[http-parent-variable][] rules raise an issue, as a significant number of
issues raised by this rule are _actually_ violations of one of those.

## Examples

### Verb only

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc CheckoutBook(CheckoutBookRequest) returns (CheckoutBookResponse) {
  option (google.api.http) = {
    // Should end with ":checkout", because the book is implied.
    post: "/v1/{name=publishers/*/books/*}:checkoutBook"
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc CheckoutBook(CheckoutBookRequest) returns (CheckoutBookResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:checkout"
    body: "*"
  };
}
```

### Verb and noun

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc SignContract(SignContractRequest) returns (SignContractResponse) {
  option (google.api.http) = {
    // Should end with ":signContract", because contract is not implied.
    post: "/v1/{publisher=publishers/*}:sign"
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc SignContract(SignContractRequest) returns (SignContractResponse) {
  option (google.api.http) = {
    post: "/v1/{publisher=publishers/*}:signContract"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0136::http-uri-suffix=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc CheckoutBook(CheckoutBookRequest) returns (CheckoutBookResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:checkoutBook"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-136]: https://aip.dev/136
[aip.dev/not-precedent]: https://aip.dev/not-precedent
[http-name-variable]: ./http-name-variable.md
[http-parent-variable]: ./http-parent-variable.md
