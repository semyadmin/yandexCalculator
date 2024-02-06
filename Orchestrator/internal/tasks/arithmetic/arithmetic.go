package arithmetic

type Expression interface {
	GetExpression() string
	Result() []string
}
