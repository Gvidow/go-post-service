package entity

type Query struct{}

type QueryConfig struct {
	Limit  int
	Cursor int
	Depth  int
}
