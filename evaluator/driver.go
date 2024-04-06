package evaluator

// Driver represents a template evaluation driver.
type Driver interface {
	// CreateContext creates a new evaluation context.
	CreateContext(*Environment) (DriverContext, error)
}

// DriverContext represents the context for template evaluation.
type DriverContext interface {
	// Eval evaluates the provided code and returns the result.
	Eval(code string) (string, error)
}
