---
rule:
  aip: 156
  name: [core, '0156', forbidden-methods]
  summary: Singletons must not define List, Create, or Delete methods.
permalink: /156/forbidden-methods
redirect_from:
  - /0156/forbidden-methods
---

# Singletons: Forbidden methods

This rule enforces that singleton resources do not define `List`, `Create`, or
`Delete` methods, as mandated in [AIP-156][].

## Details

This rule looks at any message with a `name` variable in the URI where the name
ends in anything other than `*`. It assumes that this is a method operating on
a singleton resource, and complains if the method is a `List`, `Create`, or
`Delete` standard method.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc GetSettings(GetSettingsRequest) returns (Settings) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/settings}"
  };
}

rpc UpdateSettings(UpdateSettingsRequest) returns (Settings) {
  option (google.api.http) = {
    patch: "/v1/{settings.name=publishers/*/settings}"
    body: "settings"
  };
}

// This method should not exist. The settings should be implicitly created
// when the publisher is created.
rpc CreateSettings(CreateSettingsRequest) returns (Settings) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/settings}"
    body: "settings"
  };
}

// This method should not exist. The settings should always implicitly exist.
rpc DeleteSettings(DeleteSettingsRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/settings}"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc GetSettings(GetSettingsRequest) returns (Settings) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/settings}"
  };
}

rpc UpdateSettings(UpdateSettingsRequest) returns (Settings) {
  option (google.api.http) = {
    patch: "/v1/{settings.name=publishers/*/settings}"
    body: "settings"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
rpc GetSettings(GetSettingsRequest) returns (Settings) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/settings}"
  };
}

rpc UpdateSettings(UpdateSettingsRequest) returns (Settings) {
  option (google.api.http) = {
    patch: "/v1/{settings.name=publishers/*/settings}"
    body: "settings"
  };
}

// (-- api-linter: core::0156::forbidden-methods=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc CreateSettings(CreateSettingsRequest) returns (Settings) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/settings}"
    body: "settings"
  };
}

// (-- api-linter: core::0156::forbidden-methods=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc DeleteSettings(DeleteSettingsRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    delete: "/v1/{name=publishers/*/settings}"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-156]: https://aip.dev/156
[aip.dev/not-precedent]: https://aip.dev/not-precedent
