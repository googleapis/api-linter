---
rule:
  aip: 190
  name: [core, '0190', service-case]
  summary: Service names must use UpperCamelCase.
permalink: /190/service-case
redirect_from:
  - /0190/service-case
---

# Service name case

This rule enforces that all service names use UpperCamelCase, as mandated in
[AIP-190][].

## Details

This rule checks all service names and complains if they are not UpperCamelCase,
or if they contain consecutive uppercase letters not present in the allowlist.

## Caveats

This rule is an imperfect heuristic. The root problem is that AIP-190's definition of camel case starts with the prose form of the name as input, for example "XML HTTP Request". But this prose need not actually be written down anywhere – it may only exist inside the API author's head! This linter certainly doesn't have access to it. What we actually have is just the author’s _proposed_ capitalization. While some case violations are obvious on their face, like `snake_case` or `lowerCamelCase`, the rest require guessing. 

Without knowing the actual boundaries of words, we can’t detect a violation like `XmlhttpRequest` where the letter H is incorrectly not capitalized. And even with access to a dictionary, we couldn’t tell whether a `CarpetService` should be actually capitalized as `CarPetService`! Finally, some terms can be written as a single word or multiple, like "nonempty" vs "non-empty". An API or even whole suite of APIs should probably be consistent about using one or the other but this rule isn't smart enough to know which of `CheckNotEmpty` or `CheckNonempty` is a violation. Fortunately, this undercapitalization type of error is less common and so we focus on the problem of excessive capitalization. 

As a heuristic, we flag consecutive capital letters in the proposed name, for example the `XML` part of `XMLHttpRequest`. But there are legitimate words with just a single letter, like the “x” in x-ray or the “p” in p-value. We use an allowlist of exceptions to ignore them.

## Disabling

Users should not hesitate to disable this rule if it makes a mistake. To do so, include the “prose form” of your name, as defined by the Google Java Style Guide, in a comment above the proto element. For example:

```proto
// (-- api-linter: core::0190::service-case=disabled
//     Prose service name: "T Shirt Service" --)
service TShirtService { ... }
```

[aip-190]: https://aip.dev/190
