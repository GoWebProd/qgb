package qgb

type orderBy struct {
	field string
	sort  Order
}

type Order string

const (
	Asc  Order = "ASC"
	Desc Order = "DESC"
)
