package evaluator

type HooksAPI interface {
	Include(name string) (string, error)
	Render(string, map[string]any) (string, error)
}
