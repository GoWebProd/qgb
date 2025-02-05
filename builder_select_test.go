package qgb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSelectSimple(t *testing.T) {
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
		Where(
			EQ("id"),
		).
		Build()

	assert.NoError(t, err)

	query, args := qb.Prepare(&ts)

	assert.Equal(
		t,
		`SELECT id, key, scopes, created_at, updated_at FROM "testTable" WHERE id = @id1`,
		query,
	)
	assert.Equal(t, 1, len(args))
	assert.Equal(t, &ts.ID, args["id1"])
}

func TestSelectFieldsWithLimitAndOffset(t *testing.T) {
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
		Fields("key").
		Limit(1).
		Offset(10).
		Build()

	assert.NoError(t, err)

	query, args := qb.Prepare(&ts)

	assert.Equal(
		t,
		`SELECT key FROM "testTable" LIMIT 1 OFFSET 10`,
		query,
	)
	assert.Equal(t, 0, len(args))
}

func TestSelectWithPlaceholder(t *testing.T) {
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
		Fields("key").
		Where(
			EQv("id", Placeholder("some_id")),
		).
		Limit(1).
		Offset(10).
		Build()

	assert.NoError(t, err)

	params := make(pgx.NamedArgs)
	params["some_id"] = 1234

	query, args := qb.PrepareArgs(params)

	assert.Equal(
		t,
		`SELECT key FROM "testTable" WHERE id = @some_id LIMIT 1 OFFSET 10`,
		query,
	)
	assert.Equal(t, 1, len(args))
	assert.Equal(t, params["some_id"], args["some_id"])
}

func TestSelectWithTwiceParameter(t *testing.T) {
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
		Fields("key").
		Where(
			OR(
				EQv("id", 5),
				EQv("id", 6),
			),
		).
		Build()

	assert.NoError(t, err)

	params := make(pgx.NamedArgs)
	params["id1"] = 1234
	params["id2"] = 4567

	query, args := qb.PrepareArgs(params)

	assert.Equal(
		t,
		`SELECT key FROM "testTable" WHERE (id = @id1) OR (id = @id2)`,
		query,
	)
	assert.Equal(t, 2, len(args))
	assert.Equal(t, params["id1"], args["id1"])
	assert.Equal(t, params["id2"], args["id2"])
}

func TestSelectFieldsWithLimitAndOrderBy(t *testing.T) {
	t.Parallel()

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
		Fields("key").
		OrderBy("id", Desc).
		Limit(10).
		Build()

	assert.NoError(t, err)

	query, args := qb.Prepare(&ts)

	assert.Equal(
		t,
		`SELECT key FROM "testTable" ORDER BY id DESC LIMIT 10`,
		query,
	)
	assert.Equal(t, 0, len(args))
}
