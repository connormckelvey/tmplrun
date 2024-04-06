package evaluator

// HooksAPI represents an interface for interacting with template hooks.
type HooksAPI interface {
	// Include includes the specified template.
	Include(string) (string, error)
	// Render renders the specified template with the given properties.
	Render(string, map[string]any) (string, error)
}
