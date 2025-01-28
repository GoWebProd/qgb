package qgb

import (
	"unsafe"

	"github.com/GoWebProd/gip/safe"
	"github.com/GoWebProd/gip/types/iface"
	"github.com/jackc/pgx/v5"
)

func get[T any](table *table, args []any, row pgx.Row) (*T, []any, error) {
	var t T

	if args == nil {
		args = make([]any, 0, len(table.fields)+2)
	} else {
		args = args[:0]
	}

	ptr := safe.Noescape(&t)

	for _, f := range table.fields {
		args = append(args, iface.Build(f.fType, unsafe.Add(ptr, f.offset)))
	}

	if table.createdAt != nil {
		args = append(args, iface.Build(table.createdAt.fType, unsafe.Add(ptr, table.createdAt.offset)))
	}

	if table.updatedAt != nil {
		args = append(args, iface.Build(table.updatedAt.fType, unsafe.Add(ptr, table.updatedAt.offset)))
	}

	if err := row.Scan(args...); err != nil {
		return nil, nil, err
	}

	return &t, args, nil
}

func collect[T any](table *table, row pgx.Rows) ([]*T, error) {
	var (
		args []any
		res  []*T
		err  error
	)

	defer row.Close()

	for row.Next() {
		var t *T

		t, args, err = get[T](table, args, row)
		if err != nil {
			return nil, err
		}

		res = append(res, t)
	}

	return res, nil
}
