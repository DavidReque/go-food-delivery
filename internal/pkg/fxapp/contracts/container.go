package contracts

// Container is a contract for a container that can be used to resolve dependencies.
// It is used to resolve dependencies for the application.
type Container interface {
	// ResolveFunc resolves a function with the given parameters.
	ResolveFunc(function interface{})
	// ResolveFuncWithParamTag resolves a function with the given parameters and a param tag.
	ResolveFuncWithParamTag(function interface{}, paramTagName string)
}
