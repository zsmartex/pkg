package query

type Query interface {
	Source() (interface{}, error)
}
