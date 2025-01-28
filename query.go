package qgb

import (
	"context"
	"unsafe"

	"github.com/GoWebProd/gip/fasttime"
	"github.com/GoWebProd/gip/safe"
	"github.com/GoWebProd/gip/types/iface"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Query[T any] struct {
	query  string
	table  *table
	fields []placeholderValue

	addCreatedAt string
	addUpdatedAt string
}

func (q Query[T]) String() string {
	return q.query
}

func (q Query[T]) Prepare(t *T) (string, pgx.NamedArgs) {
	args := make(pgx.NamedArgs)
	ptr := safe.Noescape(t)

	for idx := range q.fields {
		name := q.fields[idx].field
		f := q.fields[idx].value

		switch f := f.(type) {
		case *field:
			if t != nil {
				args[name] = iface.Build(f.fType, unsafe.Add(ptr, f.offset))
			}
		case placeholder:
			field, ok := q.table.fieldsMap[f.name]
			if ok {
				args[name] = iface.Build(field.fType, unsafe.Add(ptr, field.offset))
			}
		default:
			args[name] = f
		}
	}

	return q.PrepareArgs(args)
}

func (q Query[T]) PrepareArgs(args pgx.NamedArgs) (string, pgx.NamedArgs) {
	if q.addCreatedAt != "" {
		args[q.addCreatedAt] = fasttime.Now()
	}

	if q.addUpdatedAt != "" {
		args[q.addUpdatedAt] = fasttime.Now()
	}

	return q.query, args
}

func (q Query[T]) Exec(ctx context.Context, tx Querier, t *T) (int64, error) {
	query, args := q.Prepare(t)

	tag, err := tx.Exec(ctx, query, args)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}

func (q Query[T]) ExecArgs(ctx context.Context, tx Querier, args pgx.NamedArgs) (int64, error) {
	query, args := q.PrepareArgs(args)

	tag, err := tx.Exec(ctx, query, args)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}

func (q Query[T]) Query(ctx context.Context, tx Querier, t *T) (pgx.Rows, error) {
	query, args := q.Prepare(t)

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (q Query[T]) QueryStructs(ctx context.Context, tx Querier, t *T) ([]*T, error) {
	query, args := q.Prepare(t)

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	return collect[T](q.table, rows)
}

func (q Query[T]) QueryArgs(ctx context.Context, tx Querier, args pgx.NamedArgs) (pgx.Rows, error) {
	query, args := q.PrepareArgs(args)

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (q Query[T]) QueryStructsArgs(ctx context.Context, tx Querier, args pgx.NamedArgs) ([]*T, error) {
	query, args := q.PrepareArgs(args)

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	return collect[T](q.table, rows)
}

func (q Query[T]) QueryRow(ctx context.Context, tx Querier, t *T) pgx.Row {
	query, args := q.Prepare(t)

	return tx.QueryRow(ctx, query, args)
}

func (q Query[T]) QueryStruct(ctx context.Context, tx Querier, t *T) (*T, error) {
	query, args := q.Prepare(t)

	t, _, err := get[T](q.table, nil, tx.QueryRow(ctx, query, args))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (q Query[T]) QueryRowArgs(ctx context.Context, tx Querier, args pgx.NamedArgs) pgx.Row {
	query, args := q.PrepareArgs(args)

	return tx.QueryRow(ctx, query, args)
}

func (q Query[T]) QueryStructArgs(ctx context.Context, tx Querier, args pgx.NamedArgs) (*T, error) {
	query, args := q.PrepareArgs(args)

	t, _, err := get[T](q.table, nil, tx.QueryRow(ctx, query, args))
	if err != nil {
		return nil, err
	}

	return t, nil
}

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
