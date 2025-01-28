package qgb

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

type scanner struct {
	pgx.Rows

	data [][]any
	rows int
}

func (s *scanner) Close() {

}

func (s *scanner) Next() bool {
	if s.rows == 0 {
		return false
	}

	s.rows--

	return true
}

func (s *scanner) Scan(v ...any) error {
	s.data = append(s.data, v)

	return nil
}

func TestScan(t *testing.T) {
	type testStruct struct {
		ID        uint64    `db:"id,primaryKey"`
		Key       string    `db:"key"`
		Scopes    string    `db:"scopes"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	o, err := New[testStruct]("testTable")

	assert.NoError(t, err)

	scanner := &scanner{}

	ts, err := o.Get(scanner)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(scanner.data))
	assert.IsType(t, &ts.ID, scanner.data[0][0])
	assert.IsType(t, &ts.Key, scanner.data[0][1])
	assert.IsType(t, &ts.Scopes, scanner.data[0][2])
	assert.IsType(t, &ts.CreatedAt, scanner.data[0][3])
	assert.IsType(t, &ts.UpdatedAt, scanner.data[0][4])
}

func TestCollect(t *testing.T) {
	type testStruct struct {
		ID        uint64    `db:"id,primaryKey"`
		Key       string    `db:"key"`
		Scopes    string    `db:"scopes"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	o, err := New[testStruct]("testTable")

	assert.NoError(t, err)

	scanner := &scanner{rows: 3}

	ts, err := o.Collect(scanner)

	assert.NoError(t, err)

	assert.Equal(t, 3, len(scanner.data))

	for _, v := range scanner.data {
		assert.IsType(t, &ts[0].ID, v[0])
		assert.IsType(t, &ts[0].Key, v[1])
		assert.IsType(t, &ts[0].Scopes, v[2])
		assert.IsType(t, &ts[0].CreatedAt, v[3])
		assert.IsType(t, &ts[0].UpdatedAt, v[4])
	}
}
