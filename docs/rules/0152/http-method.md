---
rule:
  aip: 152
  name: [core, '0152', http-method]
  summary: Run methods must use the POST HTTP verb.
permalink: /152/http-method
redirect_from:
  - /0152/http-method
---

# Run methods: POST HTTP verb

This rule enforces that all `Run` use the `POST` HTTP verb, as
mandated in [AIP-152][].

## Details

This rule looks at any RPCs with the name beginning with `Run`, and
complains if the HTTP verb is anything other than `POST`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc RunWriteBookJob(RunWriteBookJobRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    patch: "/v1/{name=publishers/*/writeBookJobs/*}:run" // Should be `post:`.
    body: "*"
  };
  option (google.longrunning.operation_info) = {
    response_type: "RunWriteBookJobResponse"
    metadata_type: "RunWriteBookJobMetadata"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc RunWriteBookJob(RunWriteBookJobRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/writeBookJobs/*}:run"
    body: "*"
  };
  option (google.longrunning.operation_info) = {
    response_type: "RunWriteBookJobResponse"
    metadata_type: "RunWriteBookJobMetadata"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0152::http-method=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc RunWriteBookJob(RunWriteBookJobRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    patch: "/v1/{name=publishers/*/writeBookJobs/*}:run"
    body: "*"
  };
  option (google.longrunning.operation_info) = {
    response_type: "RunWriteBookJobResponse"
    metadata_type: "RunWriteBookJobMetadata"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-152]: https://aip.dev/152
[aip.dev/not-precedent]: https://aip.dev/not-precedent
