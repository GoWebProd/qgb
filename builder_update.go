package qgb

import (
	"bytes"
	"fmt"
	"strings"
)

type UpdateBuilder[T any] struct {
	table *table

	updateField []string
	updateValue []any

	where     *Clause
	returning []string

	unexpectedFields []string
}

func (b *UpdateBuilder[T]) Set(field string) *UpdateBuilder[T] {
	return b.SetValue(field, nil)
}

func (b *UpdateBuilder[T]) SetValue(field string, value any) *UpdateBuilder[T] {
	f, ok := b.table.fieldsMap[field]
	if !ok {
		b.unexpectedFields = append(b.unexpectedFields, field)

		return b
	}

	b.updateField = append(b.updateField, field)

	if value != nil {
		b.updateValue = append(b.updateValue, value)
	} else {
		b.updateValue = append(b.updateValue, f)
	}

	return b
}

func (b *UpdateBuilder[T]) Where(clause *Clause) *UpdateBuilder[T] {
	b.where = clause

	return b
}

func (b *UpdateBuilder[T]) Returning(fields ...string) *UpdateBuilder[T] {
	if fields == nil {
		fields = make([]string, 0)
	}

	b.returning = fields

	return b
}

func (b *UpdateBuilder[T]) Build() (Query[T], error) {
	var (
		q       Query[T]
		counter counter
	)

	if b.unexpectedFields != nil {
		return q, fmt.Errorf("unexpected fields: %s", strings.Join(b.unexpectedFields, ", "))
	}

	b.checkParams(&q, &counter)

	buf := bytes.NewBuffer(make([]byte, 0, 1024))

	buf.WriteString("UPDATE \"")
	buf.WriteString(b.table.name)
	buf.WriteString("\" SET ")

	for i, f := range b.updateField {
		if i != 0 {
			buf.WriteString(", ")
		}

		var pholder string

		if p, ok := b.updateValue[i].(placeholder); ok {
			pholder = p.name
		} else {
			idx := counter.IncrementString()
			pholder = f + idx
			b.updateField[i] = pholder
		}

		buf.WriteString(f)

		if f == "updated_at" {
			buf.WriteString(" = to_timestamp(@")
			buf.WriteString(pholder)
			buf.WriteString(") at time zone 'utc'")
		} else {
			buf.WriteString(" = @")
			buf.WriteString(pholder)
		}
	}

	if b.where != nil {
		sql, args, err := b.where.toSQL(&counter)
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

	for idx, f := range b.updateValue {
		if f != nil {
			name := b.updateField[idx]
			q.fields = append(q.fields, placeholderValue{field: name, value: f})
		}
	}

	if len(b.returning) > 0 {
		buf.WriteString(" RETURNING ")

		for i, f := range b.returning {
			if i != 0 {
				buf.WriteString(", ")
			}

			buf.WriteString(f)
		}
	}

	q.query = buf.String()
	q.table = b.table

	return q, nil
}

func (b *UpdateBuilder[T]) checkParams(q *Query[T], counter *counter) {
	if len(b.updateField) == 0 {
		b.updateField = make([]string, 0, len(b.table.fields)+1)
		b.updateValue = make([]any, 0, len(b.table.fields)+1)

		for _, f := range b.table.fields {
			if f.isPrimaryKey {
				continue
			}

			b.updateField = append(b.updateField, f.name)
			b.updateValue = append(b.updateValue, f)
		}
	}

	if b.table.updatedAt != nil {
		hasUpdatedAt := false

		for _, f := range b.updateField {
			if f == b.table.updatedAt.name {
				hasUpdatedAt = true

				break
			}
		}

		if !hasUpdatedAt {
			p := placeholder{"updated_at" + counter.IncrementString()}

			b.updateField = append(b.updateField, "updated_at")
			b.updateValue = append(b.updateValue, p)
			q.addUpdatedAt = p.name
		}
	}

	if b.returning != nil && len(b.returning) == 0 {
		b.returning = make([]string, 0, len(b.table.fields)+2)

		for _, f := range b.table.fields {
			b.returning = append(b.returning, f.name)
		}

		if b.table.createdAt != nil {
			b.returning = append(b.returning, "created_at")
		}

		if b.table.updatedAt != nil {
			b.returning = append(b.returning, "updated_at")
		}
	}
}
