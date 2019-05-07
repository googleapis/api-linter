# Google API Linter

Google API Linter is a tool to enforce [API Design Guide](https://cloud.google.com/apis/design/) on API Protobuf definitions by checking `.proto` files.

It differs from other Protobuf linters in that it looks at your protobuf files from a higher level and expects that the files define an API to be consumed by people not on your team (e.g., for Google Cloud Platform). For example, the API Linter will check whether your `update` methods follow the conventions defined at [API Design Guide -- Standard Methods](https://cloud.google.com/apis/design/standard_methods#update), which is out of scope for other standard Protobuf linters.
