---
rule:
  aip: 191
  name: [core, '0191', csharp-namespace]
  summary: The `option csharp_namespace` annotation should be idiomatic if set.
permalink: /191/csharp-namespace
redirect_from:
  - /0191/csharp-namespace
---

# C# namespace annotation

This rule enforces that if a proto file for a public API surface sets
`option csharp_namespace`, that it uses language idiomatic conventions, as
mandated in [AIP-191][].

## Details

This rule looks at each proto file, and complains if the `csharp_namespace`
file annotation uses anything other than upper camel case, or includes
characters other than letters, numbers, and `.`.

It also ensures that versions with stability (e.g. `V1Beta1`) are capitalized
appropriately.

## Examples

### Case

**Incorrect** code for this rule:

```proto
// Incorrect.
syntax = "proto3";

package google.example.v1;

option csharp_namespace = "google.example.v1";
```

```proto
// Incorrect.
syntax = "proto3";

package google.example.v1;

option csharp_namespace = "Google::Example::V1";
```

**Correct** code for this rule:

```proto
// Correct.
syntax = "proto3";

package google.example.v1;

option csharp_namespace = "Google.Example.V1";
```

### Versions with stability

**Incorrect** code for this rule:

```proto
// Incorrect.
syntax = "proto3";

package google.example.v1beta1;

option csharp_namespace = "Google.Example.V1beta1"; // Should be V1Beta1.
```

**Correct** code for this rule:

```proto
// Correct.
syntax = "proto3";

package google.example.v1beta1;

option csharp_namespace = "Google.Example.V1Beta1";
```

## Known issues

This rule will improperly complain if it encounters an acronym. For example, it
will complain about `Google.Cloud.AutoML.V1`, preferring `AutoMl`. This lint
rule **may** be disabled in this case.

## Disabling

If you need to violate this rule, use a comment at the top of the file.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0191::csharp-namespace=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
syntax = "proto3";

package google.example.v1;

option csharp_namespace = "google.example.v1";
```

[aip-191]: https://aip.dev/191
[aip.dev/not-precedent]: https://aip.dev/not-precedent
