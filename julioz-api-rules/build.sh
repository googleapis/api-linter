#!/bin/bash
# build.sh - Build script for the Julio API Linter Rules plugin

set -e

# Check Go version
GO_VERSION=$(go version)
echo "Building plugin with $GO_VERSION"

# This script must be run with the same version of Go as the target api-linter
# Find the installed api-linter version (if available)
API_LINTER_PATH=$(which api-linter 2>/dev/null || echo "")

if [ -n "$API_LINTER_PATH" ]; then
    INSTALLED_VERSION=$(api-linter --version 2>&1 | grep -o "api-linter [0-9.]*" | cut -d ' ' -f 2)
    echo "Found installed api-linter version: $INSTALLED_VERSION"
    # Use this version for building the plugin
    API_LINTER_VERSION="v$INSTALLED_VERSION"
else
    # Default to a specific version if api-linter is not installed
    API_LINTER_VERSION="v1.69.2"
    echo "No api-linter found in PATH. Using default version: $API_LINTER_VERSION"
    echo "Warning: For plugin compatibility, build the plugin with the same Go version as the target api-linter"
fi

# Initialize or update go.mod
if [ ! -f "go.mod" ]; then
    echo "Initializing Go module"
    go mod init github.com/julioz/julioz-api-rules
else
    echo "Using existing go.mod file"
fi

echo "Using api-linter version $API_LINTER_VERSION"

# Ensure dependencies are in sync with target api-linter
go get github.com/googleapis/api-linter@$API_LINTER_VERSION

# Tidy dependencies
go mod tidy

# Set the output file name
OUTPUT_FILE="julioz-api-rules.so"

# Build the plugin with the same Go version
echo "Building plugin..."
go build -buildmode=plugin -o $OUTPUT_FILE

# Get the absolute path of the built plugin
PLUGIN_PATH=$(realpath $OUTPUT_FILE)

echo "Plugin built: $PLUGIN_PATH"

echo ""
echo "IMPORTANT: To use this plugin, the api-linter binary must be built with:"
echo "- The exact same Go version: $GO_VERSION"
echo "- The exact same api-linter version: $API_LINTER_VERSION"
echo ""
echo "Usage: api-linter --rule-plugin=$PLUGIN_PATH your-proto-file.proto" 