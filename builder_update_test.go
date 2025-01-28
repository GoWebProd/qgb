package qgb

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestUpdateSimple(t *testing.T) {
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
		ID:     1234,
		Key:    "123",
		Scopes: "456",
	}

	qb, err := o.
		Update().
		Where(
			EQ("id"),
		).
		Returning().
		Build()

	assert.NoError(t, err)

	query, args := qb.Prepare(&ts)

	assert.Equal(
		t,
		`UPDATE "testTable" SET key = @key2, scopes = @scopes3, updated_at = to_timestamp(@updated_at1) at time zone 'utc' WHERE id = @id4 RETURNING id, key, scopes, created_at, updated_at`,
		query,
	)
	assert.Equal(t, 4, len(args), "args", args)
	assert.Equal(t, &ts.Key, args["key2"])
	assert.Equal(t, &ts.Scopes, args["scopes3"])
	assert.IsType(t, int64(0), args["updated_at1"])
	assert.Equal(t, &ts.ID, args["id4"])
}

func TestUpdateCustomSet(t *testing.T) {
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
		ID:     1234,
		Key:    "123",
		Scopes: "456",
	}
	update := "567"

	qb, err := o.
		Update().
		Set("key").
		SetValue("scopes", update).
		SetValue("id", Placeholder("test_id")).
		Where(
			EQ("id"),
		).
		Returning("id").
		Build()

	assert.NoError(t, err)

	params := make(pgx.NamedArgs)
	params["key2"] = ts.Key
	params["scopes3"] = update
	params["id4"] = ts.ID
	params["test_id"] = ts.ID

	query, args := qb.PrepareArgs(params)

	assert.Equal(
		t,
		`UPDATE "testTable" SET key = @key2, scopes = @scopes3, id = @test_id, updated_at = to_timestamp(@updated_at1) at time zone 'utc' WHERE id = @id4 RETURNING id`,
		query,
	)
	assert.Equal(t, 5, len(args), "args", args)
	assert.Equal(t, ts.Key, args["key2"])
	assert.IsType(t, update, args["scopes3"])
	assert.Equal(t, update, args["scopes3"].(string))
	assert.Equal(t, ts.ID, args["id4"])
	assert.Equal(t, ts.ID, args["test_id"])
}
