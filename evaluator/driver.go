package evaluator

type Driver interface {
	CreateContext(*Environment) (DriverContext, error)
}

type DriverContext interface {
	Eval(code string) (string, error)
}
