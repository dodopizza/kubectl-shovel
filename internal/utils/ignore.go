package utils

// Ignore invokes handler and explicitly ignores error
func Ignore(handler func() error) {
	_ = handler()
}
