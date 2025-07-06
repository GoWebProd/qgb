package qgb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClauseRaw(t *testing.T) {
	clause := RAW("id = 5")

	sql, args, err := clause.toSQL(&counter{})

	assert.NoError(t, err)
	assert.Equal(t, "id = 5", sql)
	assert.Nil(t, args)
}

func TestClauseEQ(t *testing.T) {
	clause := EQ("id")

	sql, args, err := clause.toSQL(&counter{})

	assert.NoError(t, err)
	assert.Equal(t, "id = @id1", sql)
	assert.Equal(t, []placeholderValue{{"id1", placeholder{name: "id"}}}, args)
}

func TestClauseEQv(t *testing.T) {
	clause := EQv("id", 5)

	sql, args, err := clause.toSQL(&counter{})

	assert.NoError(t, err)
	assert.Equal(t, "id = @id1", sql)
	assert.Equal(t, []placeholderValue{{"id1", 5}}, args)
}

func TestClauseNEQ(t *testing.T) {
	clause := NEQ("id")

	sql, args, err := clause.toSQL(&counter{})

	assert.NoError(t, err)
	assert.Equal(t, "id <> @id1", sql)
	assert.Equal(t, []placeholderValue{{"id1", placeholder{name: "id"}}}, args)
}

func TestClauseNEQv(t *testing.T) {
	clause := NEQv("id", 5)

	sql, args, err := clause.toSQL(&counter{})

	assert.NoError(t, err)
	assert.Equal(t, "id <> @id1", sql)
	assert.Equal(t, []placeholderValue{{"id1", 5}}, args)
}

func TestClauseGT(t *testing.T) {
	clause := GT("id")

	sql, args, err := clause.toSQL(&counter{})

	assert.NoError(t, err)
	assert.Equal(t, "id > @id1", sql)
	assert.Equal(t, []placeholderValue{{"id1", placeholder{name: "id"}}}, args)
}

func TestClauseGTv(t *testing.T) {
	clause := GTv("id", 5)

	sql, args, err := clause.toSQL(&counter{})

	assert.NoError(t, err)
	assert.Equal(t, "id > @id1", sql)
	assert.Equal(t, []placeholderValue{{"id1", 5}}, args)
}

func TestClauseGTE(t *testing.T) {
	clause := GTE("id")

	sql, args, err := clause.toSQL(&counter{})

	assert.NoError(t, err)
	assert.Equal(t, "id >= @id1", sql)
	assert.Equal(t, []placeholderValue{{"id1", placeholder{name: "id"}}}, args)
}

func TestClauseGTEv(t *testing.T) {
	clause := GTEv("id", 5)

	sql, args, err := clause.toSQL(&counter{})

	assert.NoError(t, err)
	assert.Equal(t, "id >= @id1", sql)
	assert.Equal(t, []placeholderValue{{"id1", 5}}, args)
}

func TestClauseLT(t *testing.T) {
	clause := LT("id")

	sql, args, err := clause.toSQL(&counter{})

	assert.NoError(t, err)
	assert.Equal(t, "id < @id1", sql)
	assert.Equal(t, []placeholderValue{{"id1", placeholder{name: "id"}}}, args)
}

func TestClauseLTv(t *testing.T) {
	clause := LTv("id", 5)

	sql, args, err := clause.toSQL(&counter{})

	assert.NoError(t, err)
	assert.Equal(t, "id < @id1", sql)
	assert.Equal(t, []placeholderValue{{"id1", 5}}, args)
}

func TestClauseLTE(t *testing.T) {
	clause := LTE("id")

	sql, args, err := clause.toSQL(&counter{})

	assert.NoError(t, err)
	assert.Equal(t, "id <= @id1", sql)
	assert.Equal(t, []placeholderValue{{"id1", placeholder{name: "id"}}}, args)
}

func TestClauseLTEv(t *testing.T) {
	clause := LTEv("id", 5)

	sql, args, err := clause.toSQL(&counter{})

	assert.NoError(t, err)
	assert.Equal(t, "id <= @id1", sql)
	assert.Equal(t, []placeholderValue{{"id1", 5}}, args)
}

func TestClauseIN(t *testing.T) {
	inV := []int{1, 2, 3}
	clause := IN("id", inV)

	sql, args, err := clause.toSQL(&counter{})

	assert.NoError(t, err)
	assert.Equal(t, "id IN @id1", sql)
	assert.Equal(t, []placeholderValue{{"id1", inV}}, args)
}

func TestClauseANY(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name   string
		field  string
		value  any
		assert func(t *testing.T, sql string, args []placeholderValue, err error)
	}{
		{
			name:  "ok 1 int",
			field: "id",
			value: 1,
			assert: func(t *testing.T, sql string, args []placeholderValue, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "id = ANY(@id1)", sql)
				assert.Equal(t, []placeholderValue{
					{
						field: "id1",
						value: 1,
					},
				}, args)
			},
		},
		{
			name:  "ok 2 strings",
			field: "id",
			value: []string{"1", "2"},
			assert: func(t *testing.T, sql string, args []placeholderValue, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "id = ANY(@id1)", sql)
				assert.Equal(t, []placeholderValue{
					{
						field: "id1",
						value: []string{"1", "2"},
					},
				}, args)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			clause := ANY(tc.field, tc.value)

			sql, args, err := clause.toSQL(&counter{})
			tc.assert(t, sql, args, err)
		})
	}
}

func TestClauseISNULL(t *testing.T) {
	clause := ISNULL("id")

	sql, args, err := clause.toSQL(&counter{})

	assert.NoError(t, err)
	assert.Equal(t, "id IS NULL", sql)
	assert.Equal(t, ([]placeholderValue)(nil), args)
}

func TestClauseNOTNULL(t *testing.T) {
	clause := NOTNULL("id")

	sql, args, err := clause.toSQL(&counter{})

	assert.NoError(t, err)
	assert.Equal(t, "id IS NOT NULL", sql)
	assert.Equal(t, ([]placeholderValue)(nil), args)
}

func TestClauseCONTAINS(t *testing.T) {
	clause := CONTAINS("id", 5)

	sql, args, err := clause.toSQL(&counter{})

	assert.NoError(t, err)
	assert.Equal(t, "id @> @id1", sql)
	assert.Equal(t, []placeholderValue{{"id1", 5}}, args)
}

func TestClauseAND(t *testing.T) {
	clause := AND(
		EQ("id"),
		EQ("name"),
	)

	sql, args, err := clause.toSQL(&counter{})

	assert.NoError(t, err)
	assert.Equal(t, "(id = @id1) AND (name = @name2)", sql)
	assert.Equal(t, []placeholderValue{{"id1", placeholder{name: "id"}}, {"name2", placeholder{name: "name"}}}, args)
}

func TestClauseOR(t *testing.T) {
	clause := OR(
		EQ("id"),
		EQ("name"),
	)

	sql, args, err := clause.toSQL(&counter{})

	assert.NoError(t, err)
	assert.Equal(t, "(id = @id1) OR (name = @name2)", sql)
	assert.Equal(t, []placeholderValue{{"id1", placeholder{name: "id"}}, {"name2", placeholder{name: "name"}}}, args)
}

func TestClauseNOT(t *testing.T) {
	clause := NOT(
		EQ("id"),
	)

	sql, args, err := clause.toSQL(&counter{})

	assert.NoError(t, err)
	assert.Equal(t, "NOT (id = @id1)", sql)
	assert.Equal(t, []placeholderValue{{"id1", placeholder{name: "id"}}}, args)
}

func TestClausePlaceholder(t *testing.T) {
	clause := NOT(
		EQv("id", Placeholder("test")),
	)

	sql, args, err := clause.toSQL(&counter{})

	assert.NoError(t, err)
	assert.Equal(t, "NOT (id = @test)", sql)
	assert.Equal(t, []placeholderValue{{"test", placeholder{}}}, args)
}
