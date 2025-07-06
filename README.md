# QGB - Quick Go Builder

QGB is a high-performance, type-safe SQL query builder for PostgreSQL written in Go. It provides a struct-oriented API with minimal overhead and excellent performance characteristics.

## Features

- **Type-Safe Query Building**: Uses Go generics for compile-time type safety
- **High Performance**: Outperforms popular alternatives with minimal allocations
- **PostgreSQL-Specific**: Optimized for PostgreSQL with native feature support
- **Struct-Based Mapping**: Automatic mapping between Go structs and database tables
- **Zero-Allocation Design**: Optimized memory management using unsafe operations
- **Fluent API**: Intuitive builder pattern for query construction
- **ON CONFLICT Support**: Full support for PostgreSQL's UPSERT operations
- **Automatic Timestamps**: Built-in handling for created_at and updated_at fields

## Installation

```bash
go get github.com/CookieNyanCloud/qgb
```

## Quick Start

### Define Your Model

```go
type User struct {
    ID        uint64    `db:"id,primaryKey"`
    Email     string    `db:"email"`
    Name      string    `db:"name"`
    IsActive  bool      `db:"is_active"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}
```

### Initialize ORM

```go
import (
    "github.com/GoWebProd/qgb"
    "github.com/jackc/pgx/v5"
)

// Create ORM instance for your model
orm, err := qgb.New[User]("users")
if err != nil {
    log.Fatal(err)
}
```

### Basic CRUD Operations

#### Select

```go
// Select by ID
user := &User{ID: 123}
query, err := orm.Select().
    Where(qgb.EQ("id")).
    Build()

result, err := query.QueryStruct(ctx, db, user)

// Select with multiple conditions
query, err = orm.Select().
    Where(
        qgb.AND(
            qgb.EQ("is_active"),
            qgb.GT("created_at"),
        ),
    ).
    OrderBy("created_at", qgb.DESC).
    Limit(10).
    Build()

users := []*User{}
err = query.QuerySlice(ctx, db, &users, &User{
    IsActive:  true,
    CreatedAt: time.Now().AddDate(0, -1, 0),
})
```

#### Insert

```go
// Simple insert
user := &User{
    Email:    "john@example.com",
    Name:     "John Doe",
    IsActive: true,
}

query, err := orm.Insert().
    Returning().
    Build()

createdUser, err := query.QueryStruct(ctx, db, user)

// Insert with ON CONFLICT
query, err = orm.Insert().
    OnConflict(qgb.DoUpdate("SET name = EXCLUDED.name, updated_at = NOW()", "email")).
    Returning().
    Build()

upsertedUser, err := query.QueryStruct(ctx, db, user)
```

#### Update

```go
// Update specific fields
user := &User{
    ID:       123,
    Name:     "Jane Doe",
    IsActive: false,
}

query, err := orm.Update().
    Fields("name", "is_active", "updated_at").
    Where(qgb.EQ("id")).
    Returning().
    Build()

updatedUser, err := query.QueryStruct(ctx, db, user)
```

#### Delete

```go
// Delete by ID
query, err := orm.Delete().
    Where(qgb.EQ("id")).
    Build()

_, err = query.Exec(ctx, db, &User{ID: 123})
```

## Performance Benchmarks

Benchmarked on Apple M3 Pro:

| Library | Time (ns/op) | Memory (B/op) | Allocations |
|---------|--------------|---------------|-------------|
| **QGB Prepare** | **161.6** | **344** | **3** |
| QGB Full Build | 810.1 | 1152 | 18 |
| SQLBuilder | 919.0 | 1136 | 31 |
| Bob | 2184 | 2489 | 67 |
| Squirrel | 4069 | 3722 | 79 |

QGB Prepare shows the performance when reusing prepared queries, demonstrating up to 25x better performance than some alternatives.

## Advanced Features

### WHERE Clause Operators

QGB supports a comprehensive set of operators for building WHERE clauses:

- `EQ(field)` - Equality comparison
- `NEQ(field)` - Not equal comparison
- `GT(field)` - Greater than
- `GTE(field)` - Greater than or equal
- `LT(field)` - Less than
- `LTE(field)` - Less than or equal
- `IN(field)` - IN clause for multiple values
- `ANY(field)` - PostgreSQL ANY operator
- `ISNULL(field)` - IS NULL check
- `NOTNULL(field)` - IS NOT NULL check
- `RAW(sql)` - Raw SQL clause

### Logical Operators

Combine multiple conditions using logical operators:

```go
query, err := orm.Select().
    Where(
        qgb.OR(
            qgb.AND(
                qgb.EQ("is_active"),
                qgb.GT("created_at"),
            ),
            qgb.IN("id"),
        ),
    ).
    Build()
```

### ON CONFLICT Options

PostgreSQL-specific conflict resolution:

```go
// Do nothing on conflict
query, err := orm.Insert().
    OnConflict(qgb.DoNothing("email")).
    Build()

// Update specific columns on conflict
query, err := orm.Insert().
    OnConflict(qgb.DoUpdate("SET name = EXCLUDED.name, updated_at = NOW()", "email")).
    Build()
```

### Named Parameters

QGB uses @ prefix for named parameters:

```go
// Built query will use @id, @email placeholders
query, err := orm.Select().
    Where(
        qgb.AND(
            qgb.EQ("id"),
            qgb.EQ("email"),
        ),
    ).
    Build()

// Executing with struct automatically maps fields to parameters
result, err := query.QueryStruct(ctx, db, &User{ID: 123, Email: "test@example.com"})
```

## Model Generator

QGB includes a model generator tool that can create struct definitions from existing PostgreSQL databases:

```bash
# Install the model generator
go install github.com/GoWebProd/qgb/cmd/modelgen@latest

# Generate models from your database
modelgen -dsn "postgres://user:pass@localhost/dbname" -output models/
```

## Design Philosophy

QGB is designed with the following principles:

1. **Performance First**: Every design decision prioritizes performance and minimal allocations
2. **PostgreSQL Native**: No database abstraction - fully embrace PostgreSQL features
3. **Type Safety**: Leverage Go's type system for compile-time safety
4. **Simple API**: Intuitive builder pattern that's easy to learn and use
5. **Zero Magic**: Explicit is better than implicit - no hidden behavior

## Requirements

- Go 1.23.4+
- PostgreSQL 12+
- pgx/v5 driver

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Benchmarks

To run benchmarks:

```bash
cd benchmark
go test -bench=. -benchmem
```