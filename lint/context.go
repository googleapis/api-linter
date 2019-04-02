package lint

import "context"

// Context provides additional information for a rule to perform linting.
type Context struct {
	context    context.Context
	descSource DescriptorSource
}

// NewContext creates a new `Context`.
func NewContext(ctx context.Context) Context {
	return Context{
		context: ctx,
	}
}

// NewContextWithDescriptorSource creates a new `Context` with the source.
func NewContextWithDescriptorSource(ctx context.Context, source DescriptorSource) Context {
	return Context{
		context:    ctx,
		descSource: source,
	}
}

// DescriptorSource returns a `DescriptorSource` if available; otherwise,
// returns (nil, ErrSourceInfoNotAvailable).
//
// The returned `DescriptorSource` contains additional information
// about a protobuf descriptor, such as comments and location lookups.
func (c Context) DescriptorSource() DescriptorSource {
	return c.descSource
}
