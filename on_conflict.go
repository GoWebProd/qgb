package qgb

type conflictAction string

const (
	doNothing conflictAction = "DO NOTHING"
	doUpdate  conflictAction = "DO UPDATE "
)

type OnConflict[T any] struct {
	fields []string
	action conflictAction
	set    string
	parent *InsertBuilder[T]
}

func (c *OnConflict[T]) DoNothing() *InsertBuilder[T] {
	c.action = doNothing

	return c.parent
}

func (c *OnConflict[T]) DoUpdate(set string) *InsertBuilder[T] {
	c.action = doUpdate
	c.set = set

	return c.parent
}
