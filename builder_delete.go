package qgb

import "bytes"

type DeleteBuilder[T any] struct {
	table *table

	where *clause
}

func (b *DeleteBuilder[T]) Where(clause *clause) *DeleteBuilder[T] {
	b.where = clause

	return b
}

func (b *DeleteBuilder[T]) Build() (Query[T], error) {
	var q Query[T]

	buf := bytes.NewBuffer(make([]byte, 0, 1024))

	buf.WriteString("DELETE FROM \"")
	buf.WriteString(b.table.name)
	buf.WriteString("\"")

	if b.where != nil {
		sql, args, err := b.where.toSQL(&counter{})
		if err != nil {
			return q, err
		}

		buf.WriteString(" WHERE ")
		buf.WriteString(sql)

		if args != nil {
			q.fields, err = transformArgs(b.table, args)
			if err != nil {
				return q, err
			}
		}
	}

	q.query = buf.String()
	q.table = b.table

	return q, nil
}
