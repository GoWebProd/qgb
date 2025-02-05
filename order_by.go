package qgb

type orderBy struct {
	field string
	sort  sort
}

type sort string

const (
	Asc  sort = "ASC"
	Desc sort = "DESC"
)
