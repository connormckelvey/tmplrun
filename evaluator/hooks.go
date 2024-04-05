package evaluator

type HooksAPI interface {
	Include(string) (string, error)
	Render(string, map[string]any) (string, error)
}
