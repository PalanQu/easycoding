package orm

type Model interface {
	TableName() string
}
