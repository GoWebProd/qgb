package qgb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDeleteSimple(t *testing.T) {
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

	qb, err := o.Delete().Where(EQ("id")).Build()

	assert.NoError(t, err)

	query, args := qb.Prepare(&ts)

	assert.Equal(
		t,
		`DELETE FROM "testTable" WHERE id = @id1`,
		query,
	)
	assert.Equal(t, 1, len(args), "args %v", args)
	assert.Equal(t, &ts.ID, args["id1"])
}
