package qgb

import (
	"bytes"
	"strconv"
	"strings"
)

type SelectBuilder[T any] struct {
	table *table

	fields       []string
	fieldsCustom bool

	where   *clause
	orderBy *orderBy
	limit   *int
	offset  *int
}

func (b *SelectBuilder[T]) Fields(fields ...string) *SelectBuilder[T] {
	if fields != nil {
		b.fieldsCustom = true
	}

	b.fields = fields

	return b
}

func (b *SelectBuilder[T]) Where(clause *clause) *SelectBuilder[T] {
	b.where = clause

	return b
}

func (b *SelectBuilder[T]) Limit(limit int) *SelectBuilder[T] {
	b.limit = &limit

	return b
}

func (b *SelectBuilder[T]) Offset(offset int) *SelectBuilder[T] {
	b.offset = &offset

	return b
}

func (b *SelectBuilder[T]) OrderBy(filed string, sort sort) *SelectBuilder[T] {
	b.orderBy = &orderBy{
		field: filed,
		sort:  sort,
	}

	return b
}

func (b *SelectBuilder[T]) Build() (Query[T], error) {
	var q Query[T]

	b.checkParams()

	buf := bytes.NewBuffer(make([]byte, 0, 1024))

	buf.WriteString("SELECT ")
	buf.WriteString(strings.Join(b.fields, ", "))
	buf.WriteString(" FROM \"")
	buf.WriteString(b.table.name)
	buf.WriteString("\"")

	if b.where != nil {
		sql, args, err := b.where.toSQL(&counter{})
		if err != nil {
			return q, err
		}

		if args != nil {
			q.fields, err = transformArgs(b.table, args)
			if err != nil {
				return q, err
			}
		}

		buf.WriteString(" WHERE ")
		buf.WriteString(sql)
	}

	if b.orderBy != nil {
		buf.WriteString(" ORDER BY ")
		buf.WriteString(b.orderBy.field)
		buf.WriteString(" ")
		buf.WriteString(string(b.orderBy.sort))
	}

	if b.limit != nil {
		buf.WriteString(" LIMIT ")
		buf.WriteString(strconv.Itoa(*b.limit))
	}

	if b.offset != nil {
		buf.WriteString(" OFFSET ")
		buf.WriteString(strconv.Itoa(*b.offset))
	}

	q.query = buf.String()
	q.table = b.table

	return q, nil
}

func (b *SelectBuilder[T]) checkParams() {
	if len(b.fields) != 0 {
		return
	}

	b.fields = make([]string, len(b.table.fields))

	for i, f := range b.table.fields {
		b.fields[i] = f.name
	}

	if b.table.createdAt != nil {
		b.fields = append(b.fields, "created_at")
	}

	if b.table.updatedAt != nil {
		b.fields = append(b.fields, "updated_at")
	}
}
