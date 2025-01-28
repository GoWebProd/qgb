package qgb

import (
	"bytes"
	"fmt"
	"strings"
)

type InsertBuilder[T any] struct {
	table *table

	fields         []string
	skipPrimaryKey bool

	returning       []string
	returningCustom bool
}

func (b *InsertBuilder[T]) Fields(fields ...string) *InsertBuilder[T] {
	b.fields = fields

	return b
}

func (b *InsertBuilder[T]) Returning(fields ...string) *InsertBuilder[T] {
	if fields == nil {
		fields = make([]string, 0)
	} else {
		b.returningCustom = true
	}

	b.returning = fields

	return b
}

func (b *InsertBuilder[T]) SkipPrimaryKey() *InsertBuilder[T] {
	b.skipPrimaryKey = true

	return b
}

func (b *InsertBuilder[T]) Build() (Query[T], error) {
	var q Query[T]

	q.fields = make([]placeholderValue, 0, len(b.table.fields))

	b.checkParams()

	insertFields := make([]string, 0, len(b.table.fields))
	valuesFields := make([]string, 0, len(b.table.fields))
	returnFields := make([]string, 0, len(b.table.fields))

	for _, f := range b.fields {
		field, ok := b.table.fieldsMap[f]
		if !ok {
			return q, fmt.Errorf("field %s not found in table %s", f, b.table.name)
		}

		if field.isPrimaryKey && b.skipPrimaryKey {
			continue
		}

		insertFields = append(insertFields, f)
		valuesFields = append(valuesFields, "@"+f)
		q.fields = append(q.fields, placeholderValue{field: field.name, value: field})
	}

	for _, f := range b.returning {
		if _, ok := b.table.fieldsMap[f]; !ok {
			return q, fmt.Errorf("field %s not found in table %s", f, b.table.name)
		}

		returnFields = append(returnFields, f)
	}

	if b.table.createdAt != nil {
		insertFields = append(insertFields, "created_at")
		valuesFields = append(valuesFields, "to_timestamp(@created_at) at time zone 'utc'")
		q.addCreatedAt = "created_at"

		if b.returning != nil && !b.returningCustom {
			returnFields = append(returnFields, "created_at")
		}
	}

	if b.table.updatedAt != nil {
		insertFields = append(insertFields, "updated_at")
		valuesFields = append(valuesFields, "to_timestamp(@updated_at) at time zone 'utc'")
		q.addUpdatedAt = "updated_at"

		if b.returning != nil && !b.returningCustom {
			returnFields = append(returnFields, "updated_at")
		}
	}

	buf := bytes.NewBuffer(make([]byte, 0, 1024))

	buf.WriteString("INSERT INTO \"")
	buf.WriteString(b.table.name)
	buf.WriteString("\" (")
	buf.WriteString(strings.Join(insertFields, ", "))
	buf.WriteString(") VALUES (")
	buf.WriteString(strings.Join(valuesFields, ", "))
	buf.WriteString(")")

	if b.returning != nil {
		buf.WriteString(" RETURNING ")
		buf.WriteString(strings.Join(returnFields, ", "))
	}

	q.query = buf.String()
	q.table = b.table

	return q, nil
}

func (b *InsertBuilder[T]) checkParams() {
	if len(b.fields) == 0 {
		b.fields = make([]string, len(b.table.fields))

		for i, f := range b.table.fields {
			b.fields[i] = f.name
		}
	}

	if b.returning != nil && len(b.returning) == 0 {
		b.returning = make([]string, len(b.table.fields))

		for i, f := range b.table.fields {
			b.returning[i] = f.name
		}
	}
}
