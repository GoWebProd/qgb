package qgb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInsertSkipAndReturning(t *testing.T) {
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
		Key:    "123",
		Scopes: "456",
	}

	qb, err := o.Insert().SkipPrimaryKey().Returning().Build()
	assert.NoError(t, err)

	query, args := qb.Prepare(&ts)

	assert.Equal(
		t,
		`INSERT INTO "testTable" (key, scopes, created_at, updated_at) VALUES (@key, @scopes, to_timestamp(@created_at) at time zone 'utc', to_timestamp(@updated_at) at time zone 'utc') RETURNING id, key, scopes, created_at, updated_at`,
		query,
	)
	assert.Equal(t, 4, len(args))
	assert.Equal(t, &ts.Key, args["key"])
	assert.Equal(t, &ts.Scopes, args["scopes"])
	assert.IsType(t, int64(0), args["created_at"])
	assert.IsType(t, int64(0), args["updated_at"])
}

func TestInsertNoSkipAndNoReturning(t *testing.T) {
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
		ID:     5,
		Key:    "123",
		Scopes: "456",
	}

	qb, err := o.Insert().Build()
	assert.NoError(t, err)

	query, args := qb.Prepare(&ts)

	assert.Equal(
		t,
		`INSERT INTO "testTable" (id, key, scopes, created_at, updated_at) VALUES (@id, @key, @scopes, to_timestamp(@created_at) at time zone 'utc', to_timestamp(@updated_at) at time zone 'utc')`,
		query,
	)
	assert.Equal(t, 5, len(args))
	assert.Equal(t, &ts.ID, args["id"])
	assert.Equal(t, &ts.Key, args["key"])
	assert.Equal(t, &ts.Scopes, args["scopes"])
	assert.IsType(t, int64(0), args["created_at"])
	assert.IsType(t, int64(0), args["updated_at"])
}

func TestInsertFieldsAndCustomReturning(t *testing.T) {
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
		Key:    "123",
		Scopes: "456",
	}

	qb, err := o.Insert().Fields("key").Returning("id").Build()
	assert.NoError(t, err)

	query, args := qb.Prepare(&ts)

	assert.Equal(
		t,
		`INSERT INTO "testTable" (key, created_at, updated_at) VALUES (@key, to_timestamp(@created_at) at time zone 'utc', to_timestamp(@updated_at) at time zone 'utc') RETURNING id`,
		query,
	)
	assert.Equal(t, 3, len(args))
	assert.Equal(t, &ts.Key, args["key"])
	assert.IsType(t, int64(0), args["created_at"])
	assert.IsType(t, int64(0), args["updated_at"])
}

func TestInsertOnConflictDoNothing(t *testing.T) {
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
		Key:    "123",
		Scopes: "456",
	}

	qb, err := o.
		Insert().
		Fields("key").
		OnConflict().DoNothing().
		Returning("id").
		Build()
	assert.NoError(t, err)

	query, args := qb.Prepare(&ts)

	assert.Equal(
		t,
		`INSERT INTO "testTable" (key, created_at, updated_at) VALUES (@key, to_timestamp(@created_at) at time zone 'utc', to_timestamp(@updated_at) at time zone 'utc') ON CONFLICT DO NOTHING RETURNING id`,
		query,
	)
	assert.Equal(t, 3, len(args))
	assert.Equal(t, &ts.Key, args["key"])
	assert.IsType(t, int64(0), args["created_at"])
	assert.IsType(t, int64(0), args["updated_at"])
}

func TestInsertOnConflictDoUpdate(t *testing.T) {
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
		Key:    "123",
		Scopes: "456",
	}

	qb, err := o.
		Insert().
		Fields("key").
		OnConflict("created_at", "updated_at").DoUpdate("set created_at=testTable.created_at").
		Returning("id").
		Build()
	assert.NoError(t, err)

	query, args := qb.Prepare(&ts)

	assert.Equal(
		t,
		`INSERT INTO "testTable" (key, created_at, updated_at) VALUES (@key, to_timestamp(@created_at) at time zone 'utc', to_timestamp(@updated_at) at time zone 'utc') ON CONFLICT (created_at, updated_at) DO UPDATE set created_at=testTable.created_at RETURNING id`,
		query,
	)
	assert.Equal(t, 3, len(args))
	assert.Equal(t, &ts.Key, args["key"])
	assert.IsType(t, int64(0), args["created_at"])
	assert.IsType(t, int64(0), args["updated_at"])
}
