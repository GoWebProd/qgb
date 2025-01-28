package qgb

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

type executor struct {
	pgx.Tx

	expectedQuery string
	expectedArgs  []any
	t             *testing.T
	scanner       scanner
}

func (e *executor) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	assert.Equal(e.t, e.expectedQuery, sql)
	assert.Equal(e.t, e.expectedArgs, args)

	return pgconn.NewCommandTag("SELECT 100"), nil
}

func (e *executor) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	assert.Equal(e.t, e.expectedQuery, sql)
	assert.Equal(e.t, e.expectedArgs, args)

	return &e.scanner, nil
}

func (e *executor) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	assert.Equal(e.t, e.expectedQuery, sql)
	assert.Equal(e.t, e.expectedArgs, args)

	return &e.scanner
}

func TestQueryExec(t *testing.T) {
	type testStruct struct {
		ID        uint64    `db:"id,primaryKey"`
		Key       string    `db:"key"`
		Scopes    string    `db:"scopes"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	o, err := New[testStruct]("testTable")

	assert.NoError(t, err)

	ts := testStruct{
		ID: 1234,
	}

	qb, err := o.
		Select().
		Build()

	assert.NoError(t, err)

	executor := &executor{
		t:             t,
		expectedQuery: `SELECT id, key, scopes, created_at, updated_at FROM "testTable"`,
		expectedArgs:  []any{pgx.NamedArgs{}},
	}

	_, err = qb.Exec(context.Background(), executor, &ts)

	assert.NoError(t, err)
}

func TestQueryExecArgs(t *testing.T) {
	type testStruct struct {
		ID        uint64    `db:"id,primaryKey"`
		Key       string    `db:"key"`
		Scopes    string    `db:"scopes"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	o, err := New[testStruct]("testTable")

	assert.NoError(t, err)

	qb, err := o.
		Select().
		Build()

	assert.NoError(t, err)

	args := make(pgx.NamedArgs)

	executor := &executor{
		t:             t,
		expectedQuery: `SELECT id, key, scopes, created_at, updated_at FROM "testTable"`,
		expectedArgs:  []any{args},
	}

	rows, err := qb.ExecArgs(context.Background(), executor, args)

	assert.NoError(t, err)
	assert.Equal(t, int64(100), rows)
}

func TestQueryQuery(t *testing.T) {
	type testStruct struct {
		ID        uint64    `db:"id,primaryKey"`
		Key       string    `db:"key"`
		Scopes    string    `db:"scopes"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	o, err := New[testStruct]("testTable")

	assert.NoError(t, err)

	ts := testStruct{
		ID: 1234,
	}

	qb, err := o.
		Select().
		Build()

	assert.NoError(t, err)

	executor := &executor{
		t:             t,
		expectedQuery: `SELECT id, key, scopes, created_at, updated_at FROM "testTable"`,
		expectedArgs:  []any{pgx.NamedArgs{}},
	}

	_, err = qb.Query(context.Background(), executor, &ts)

	assert.NoError(t, err)
}

func TestQueryQueryStructs(t *testing.T) {
	type testStruct struct {
		ID        uint64    `db:"id,primaryKey"`
		Key       string    `db:"key"`
		Scopes    string    `db:"scopes"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	o, err := New[testStruct]("testTable")

	assert.NoError(t, err)

	qb, err := o.
		Select().
		Build()

	assert.NoError(t, err)

	ts := testStruct{
		ID: 1234,
	}

	executor := &executor{
		t:             t,
		expectedQuery: `SELECT id, key, scopes, created_at, updated_at FROM "testTable"`,
		expectedArgs:  []any{pgx.NamedArgs{}},
		scanner:       scanner{rows: 5},
	}

	rows, err := qb.QueryStructs(context.Background(), executor, &ts)

	assert.NoError(t, err)
	assert.Equal(t, 5, len(rows))

	for _, v := range executor.scanner.data {
		assert.IsType(t, &rows[0].ID, v[0])
		assert.IsType(t, &rows[0].Key, v[1])
		assert.IsType(t, &rows[0].Scopes, v[2])
		assert.IsType(t, &rows[0].CreatedAt, v[3])
		assert.IsType(t, &rows[0].UpdatedAt, v[4])
	}
}

func TestQueryQueryArgs(t *testing.T) {
	type testStruct struct {
		ID        uint64    `db:"id,primaryKey"`
		Key       string    `db:"key"`
		Scopes    string    `db:"scopes"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	o, err := New[testStruct]("testTable")

	assert.NoError(t, err)

	qb, err := o.
		Select().
		Build()

	assert.NoError(t, err)

	args := make(pgx.NamedArgs)

	executor := &executor{
		t:             t,
		expectedQuery: `SELECT id, key, scopes, created_at, updated_at FROM "testTable"`,
		expectedArgs:  []any{args},
		scanner:       scanner{rows: 5},
	}

	_, err = qb.QueryArgs(context.Background(), executor, args)

	assert.NoError(t, err)
}

func TestQueryQueryStructsArgs(t *testing.T) {
	type testStruct struct {
		ID        uint64    `db:"id,primaryKey"`
		Key       string    `db:"key"`
		Scopes    string    `db:"scopes"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	o, err := New[testStruct]("testTable")

	assert.NoError(t, err)

	qb, err := o.
		Select().
		Build()

	assert.NoError(t, err)

	args := pgx.NamedArgs{}

	executor := &executor{
		t:             t,
		expectedQuery: `SELECT id, key, scopes, created_at, updated_at FROM "testTable"`,
		expectedArgs:  []any{args},
		scanner:       scanner{rows: 5},
	}

	rows, err := qb.QueryStructsArgs(context.Background(), executor, args)

	assert.NoError(t, err)
	assert.Equal(t, 5, len(rows))

	for _, v := range executor.scanner.data {
		assert.IsType(t, &rows[0].ID, v[0])
		assert.IsType(t, &rows[0].Key, v[1])
		assert.IsType(t, &rows[0].Scopes, v[2])
		assert.IsType(t, &rows[0].CreatedAt, v[3])
		assert.IsType(t, &rows[0].UpdatedAt, v[4])
	}
}

func TestQueryQueryRow(t *testing.T) {
	type testStruct struct {
		ID        uint64    `db:"id,primaryKey"`
		Key       string    `db:"key"`
		Scopes    string    `db:"scopes"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	o, err := New[testStruct]("testTable")

	assert.NoError(t, err)

	ts := testStruct{
		ID: 1234,
	}

	qb, err := o.
		Select().
		Build()

	assert.NoError(t, err)

	executor := &executor{
		t:             t,
		expectedQuery: `SELECT id, key, scopes, created_at, updated_at FROM "testTable"`,
		expectedArgs:  []any{pgx.NamedArgs{}},
	}

	row := qb.QueryRow(context.Background(), executor, &ts)

	assert.NoError(t, row.Scan())
}

func TestQueryQueryStruct(t *testing.T) {
	type testStruct struct {
		ID        uint64    `db:"id,primaryKey"`
		Key       string    `db:"key"`
		Scopes    string    `db:"scopes"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	o, err := New[testStruct]("testTable")

	assert.NoError(t, err)

	qb, err := o.
		Select().
		Build()

	assert.NoError(t, err)

	ts := testStruct{
		ID: 1234,
	}

	executor := &executor{
		t:             t,
		expectedQuery: `SELECT id, key, scopes, created_at, updated_at FROM "testTable"`,
		expectedArgs:  []any{pgx.NamedArgs{}},
		scanner:       scanner{rows: 5},
	}

	row, err := qb.QueryStruct(context.Background(), executor, &ts)

	assert.NoError(t, err)

	assert.IsType(t, &row.ID, executor.scanner.data[0][0])
	assert.IsType(t, &row.Key, executor.scanner.data[0][1])
	assert.IsType(t, &row.Scopes, executor.scanner.data[0][2])
	assert.IsType(t, &row.CreatedAt, executor.scanner.data[0][3])
	assert.IsType(t, &row.UpdatedAt, executor.scanner.data[0][4])
}

func TestQueryQueryRowArgs(t *testing.T) {
	type testStruct struct {
		ID        uint64    `db:"id,primaryKey"`
		Key       string    `db:"key"`
		Scopes    string    `db:"scopes"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	o, err := New[testStruct]("testTable")

	assert.NoError(t, err)

	qb, err := o.
		Select().
		Build()

	assert.NoError(t, err)

	args := make(pgx.NamedArgs)

	executor := &executor{
		t:             t,
		expectedQuery: `SELECT id, key, scopes, created_at, updated_at FROM "testTable"`,
		expectedArgs:  []any{args},
	}

	row := qb.QueryRowArgs(context.Background(), executor, args)

	assert.NoError(t, row.Scan())
}

func TestQueryQueryStructArgs(t *testing.T) {
	type testStruct struct {
		ID        uint64    `db:"id,primaryKey"`
		Key       string    `db:"key"`
		Scopes    string    `db:"scopes"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	o, err := New[testStruct]("testTable")

	assert.NoError(t, err)

	qb, err := o.
		Select().
		Build()

	assert.NoError(t, err)

	args := pgx.NamedArgs{}

	executor := &executor{
		t:             t,
		expectedQuery: `SELECT id, key, scopes, created_at, updated_at FROM "testTable"`,
		expectedArgs:  []any{args},
		scanner:       scanner{rows: 1},
	}

	row, err := qb.QueryStructArgs(context.Background(), executor, args)

	assert.NoError(t, err)

	assert.IsType(t, &row.ID, executor.scanner.data[0][0])
	assert.IsType(t, &row.Key, executor.scanner.data[0][1])
	assert.IsType(t, &row.Scopes, executor.scanner.data[0][2])
	assert.IsType(t, &row.CreatedAt, executor.scanner.data[0][3])
	assert.IsType(t, &row.UpdatedAt, executor.scanner.data[0][4])
}
