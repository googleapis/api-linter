# Julio's API Linter Rules

This plugin provides custom rules for the [Google API Linter](https://github.com/googleapis/api-linter).

## Rules

- `core::9000::julio-prefix`: Message names must start with "Julio"

## Building the Plugin

1. Make sure you have the same Go version installed as the api-linter you're targeting
2. Run the build script:

```bash
chmod +x build.sh
./build.sh
```

This will produce `julioz-api-rules.so` that can be used with api-linter.

## Using the Plugin

Once the api-linter has plugin support, you'll be able to use this plugin with:

```bash
api-linter --rule-plugin=/path/to/julioz-api-rules.so your-proto-files
```

## Compatibility Notes

This plugin must be built with the **exact same version** of Go as the api-linter binary you're using. Additionally, all shared dependencies must match exactly. The build script tries to ensure this, but you may need to adjust the API_LINTER_VERSION variable in `build.sh` to match your installation.

## Development

The plugin implements a single rule that enforces message names to start with "Julio". To add more rules, modify the `AddCustomRules` function in `main.go`.

Unit tests are available in `rule_test.go` and can be run with:

```bash
go test ./...
```
