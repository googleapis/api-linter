package pragma

// DoNotImplement can be embedded in an interface to prevent trivial
// implementations of the interface.
//
// This is useful to prevent unauthorized implementations of an interface
// so that it can be extended in the future for any API changes.
type DoNotImplement interface{ LinterInternal(DoNotImplement) }
