---
---

# Google API Linter

The API linter provides real-time checks for compliance with many of Google's
API standards, documented using [API Improvement Proposals][]. It operates on
API surfaces defined in [protocol buffers][].

It identifies common mistakes and inconsistencies in API surfaces:

```proto
// Incorrect.
message GetBookRequest {
  // This is wrong; it should be spelled `name`.
  string book = 1;
}
```

When able, it also offers a suggestion for the correct fix.

**Note:** Not every piece of AIP guidance is able to be expressed as lint rules
(and some things that are able to be expressed may not be written yet). The
linter should be used as a useful tool, but not as a substitute for reading and
understanding API guidance.

Each linter rule has its own [rule documentation][], and rules can be
[configured][configuration] using a config file, or in a proto file itself.

## Installation

Lorem ipsum dolor set amet...

## License

This software is made available under the [Apache 2.0][] license.

[apache 2.0]: https://www.apache.org/licenses/LICENSE-2.0
[api improvement proposals]: https://aip.dev/
[protocol buffers]: https://developers.google.com/protocol-buffers
[rule documentation]: ./rules/index.md
