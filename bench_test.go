package qgb_test

import (
	"context"
	"testing"
	"time"

	"github.com/GoWebProd/qgb"

	sq "github.com/Masterminds/squirrel"
	sqb "github.com/huandu/go-sqlbuilder"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/um"
	"github.com/stretchr/testify/assert"
)

func BenchmarkFullBuild(b *testing.B) {
	type testStruct struct {
		ID        uint64    `db:"id,primaryKey"`
		Key       string    `db:"key"`
		Scopes    string    `db:"scopes"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	o, err := qgb.New[testStruct]("testTable")

	assert.NoError(b, err)

	ts := testStruct{
		ID:     1234,
		Key:    "123",
		Scopes: "456",
	}

	for i := 0; i < b.N; i++ {
		qb, _ := o.
			Update().
			Where(
				qgb.EQ("id"),
			).
			Returning().
			Build()

		qb.Prepare(&ts)
	}
}

func BenchmarkPrepare(b *testing.B) {
	type testStruct struct {
		ID        uint64    `db:"id,primaryKey"`
		Key       string    `db:"key"`
		Scopes    string    `db:"scopes"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	o, err := qgb.New[testStruct]("testTable")

	assert.NoError(b, err)

	ts := testStruct{
		ID:     1234,
		Key:    "123",
		Scopes: "456",
	}

	qb, _ := o.
		Update().
		Where(
			qgb.EQ("id"),
		).
		Returning().
		Build()

	for i := 0; i < b.N; i++ {

		qb.Prepare(&ts)
	}
}

func BenchmarkSquirrel(b *testing.B) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	for i := 0; i < b.N; i++ {
		psql.
			Update("testTable").
			Set("key", "123").
			Set("scopes", "456").
			Where("id = ?", 1234).
			Suffix("RETURNING id, key, scopes, created_at, updated_at").
			ToSql()
	}
}

func BenchmarkSQLBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sb := sqb.NewUpdateBuilder()
		sb.
			Update("testTable").
			SetMore("key", "123").
			SetMore("scopes", "456").
			Where(sb.EQ("id", 1234)).
			Build()
	}
}

func BenchmarkBob(b *testing.B) {
	for i := 0; i < b.N; i++ {
		psql.Update(
			um.Table("users"),
			um.SetCol("key").ToArg("123"),
			um.SetCol("scopes").ToArg("456"),
			um.Where(psql.Quote("id").EQ(psql.Arg(1234))),
		).Build(context.Background())
	}
}
