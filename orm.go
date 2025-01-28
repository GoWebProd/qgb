package qgb

import (
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type ORM[T any] struct {
	table table
}

func New[T any](tableName string) (*ORM[T], error) {
	var t T

	if reflect.TypeOf(t).Kind() != reflect.Struct {
		return nil, errors.Errorf("bad type passed to handler: %T, need struct", t)
	}

	var (
		orm ORM[T]
		err error
	)

	orm.table, err = buildTable[T](tableName)
	if err != nil {
		return nil, err
	}

	return &orm, nil
}

func (o *ORM[T]) Get(row pgx.Row) (*T, error) {
	t, _, err := get[T](&o.table, nil, row)

	return t, err
}

func (o *ORM[T]) Collect(row pgx.Rows) ([]*T, error) {
	return collect[T](&o.table, row)
}

func (o *ORM[T]) Insert() *InsertBuilder[T] {
	return &InsertBuilder[T]{
		table: &o.table,
	}
}

func (o *ORM[T]) Select() *SelectBuilder[T] {
	return &SelectBuilder[T]{
		table: &o.table,
	}
}

func (o *ORM[T]) Update() *UpdateBuilder[T] {
	return &UpdateBuilder[T]{
		table: &o.table,
	}
}

func (o *ORM[T]) Delete() *DeleteBuilder[T] {
	return &DeleteBuilder[T]{
		table: &o.table,
	}
}
