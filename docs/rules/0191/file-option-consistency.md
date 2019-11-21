---
rule:
  aip: 191
  name: [core, '0191', file-option-consistency]
  summary: All proto files must set file packaging options consistently.
permalink: /191/file-option-consistency
redirect_from:
  - /0191/file-option-consistency
---

# Java package annotation

This rule enforces that every proto file for a public API surface sets file
packaging options consistently, as mandated in [AIP-191][].

## Details

This rule looks at each proto file, and reads any files that it imports that
are in the same proto package. It iterates over the file packaging options in
each one and complains if they are inconsistent.

The following annotations are included:

- `csharp_namespace`
- `go_package`
- `java_multiple_files`
- `java_package`
- `php_class_prefix`
- `php_metadata_namespace`
- `php_namespace`
- `objc_class_prefix`
- `ruby_package`
- `swift_prefix`

## Examples

**Incorrect** code for this rule:

In `foo.proto`:

```proto
// Incorrect.
syntax = "proto3";

package google.example.v1;

option csharp_namespace = "Google\\Example\\V1";
option ruby_namespace = "Google::Example::V1";
```

In `bar.proto`:

```proto
// Incorrect.
syntax = "proto3";

package google.example.v1;

import "foo.proto";

option csharp_namespace = "Example\\V1";  // Inconsistent.
// option ruby_namespace is missing, which is also inconsistent.
```

**Correct** code for this rule:

In `foo.proto`:

```proto
// Correct.
syntax = "proto3";

package google.example.v1;

option csharp_namespace = "Google\\Example\\V1";
option ruby_namespace = "Google::Example::V1";
```

In `bar.proto`:

```proto
// Correct.
syntax = "proto3";

package google.example.v1;

import "foo.proto";

option csharp_namespace = "Google\\Example\\V1";
option ruby_namespace = "Google::Example::V1";
```

## Disabling

If you need to violate this rule, use a comment at the top of the file.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0191::file-option-consistency=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
syntax = "proto3";

package google.example.v1;

import "foo.proto";

option csharp_namespace = "Example\\V1";
```

[aip-191]: https://aip.dev/191
[aip.dev/not-precedent]: https://aip.dev/not-precedent
