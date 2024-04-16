package evaluator

// HooksAPI represents an interface for interacting with template hooks.
type HooksAPI interface {
	// Include includes the specified template.
	Include(string) (string, error)
	// Render renders the specified template with the given properties.
	Render(string, map[string]any) (string, error)
}

// RegisterHooks represents an interface for registering custom hooks
type RegisterHooks interface {
	Register(onError func(error), reg func(name string, value any) error) error
}
