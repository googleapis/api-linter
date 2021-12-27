---
rule:
  aip: 193
  name: [core, '0193', response-status-type]
  summary: Services must return a google.rpc.Status message when an API error occurs
permalink: /193/response-status-type
redirect_from:
  - /0193/response-status-type
---

# Response status type

This rule attempts to enforce that every response message uses the same
standard type for status objects, as mandated in [AIP-193][].

## Details

This rule looks at each response message type of every service in proto file and
verifies whether messages with the names like `status`, `error`, or `warning`
have the correct type (`google.rpc.Status`). Only fields in the top-level
message are checked.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message GetBookResponse {
  int32 status = 1;
}
```

**Correct** code for this rule:

```proto
// Correct.
import "google/rpc/status.proto";

message GetBookResponse {
  google.rpc.Status status = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the descriptor.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0193::response-status-type=disabled
//     aip.dev/not-precedent: Required for backwards compability reasons. --)
message GetBookStatus {
  int32 status = 1;
}
```

[aip-193]: https://aip.dev/193
[aip.dev/not-precedent]: https://aip.dev/not-precedent
