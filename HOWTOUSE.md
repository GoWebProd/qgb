# QGB - How To Use Guide for AI Agents

This guide provides comprehensive instructions for AI agents and developers on how to effectively use the QGB (Quick Go Builder) library.

## Table of Contents
1. [Getting Started](#getting-started)
2. [Core Concepts](#core-concepts)
3. [Query Building Patterns](#query-building-patterns)
4. [Advanced Features](#advanced-features)
5. [Performance Best Practices](#performance-best-practices)
6. [Common Patterns](#common-patterns)
7. [Error Handling](#error-handling)
8. [Testing Strategies](#testing-strategies)

## Getting Started

### Prerequisites
- Go 1.23.4+ (generics required)
- PostgreSQL 12+
- pgx/v5 driver

### Installation
```bash
go get github.com/GoWebProd/qgb
```

### Basic Setup Pattern
```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/GoWebProd/qgb"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

// Step 1: Define your model with proper db tags
type User struct {
    ID        uint64    `db:"id,primaryKey"`
    Email     string    `db:"email"`
    Name      string    `db:"name"`
    IsActive  bool      `db:"is_active"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

// Step 2: Initialize ORM
func setupORM() (*qgb.ORM[User], error) {
    return qgb.New[User]("users")
}

// Step 3: Setup database connection
func setupDB() (*pgxpool.Pool, error) {
    config, err := pgxpool.ParseConfig("postgres://user:pass@localhost/dbname")
    if err != nil {
        return nil, err
    }
    
    return pgxpool.NewWithConfig(context.Background(), config)
}
```

## Core Concepts

### 1. Struct Tags
QGB uses `db` tags to map Go struct fields to database columns:

```go
type Product struct {
    ID          uint64    `db:"id,primaryKey"`           // Primary key
    Name        string    `db:"name"`                     // Regular field
    Price       float64   `db:"price"`                    // Numeric field
    CategoryID  *uint64   `db:"category_id"`              // Nullable foreign key
    IsActive    bool      `db:"is_active"`                // Boolean field
    CreatedAt   time.Time `db:"created_at"`               // Timestamp
    UpdatedAt   time.Time `db:"updated_at"`               // Auto-updated timestamp
}
```

### 2. ORM Initialization
Always initialize the ORM with the correct table name:

```go
// Generic type must match your struct
orm, err := qgb.New[Product]("products")
if err != nil {
    return fmt.Errorf("failed to initialize ORM: %w", err)
}
```

### 3. Query Building Pattern
All queries follow the same pattern:
1. Start with a builder method (Select, Insert, Update, Delete)
2. Chain configuration methods
3. Call Build() to create the query
4. Execute with appropriate Query method

```go
// Pattern: orm.Operation().Config().Build()
query, err := orm.Select().
    Where(qgb.EQ("id")).
    OrderBy("created_at", qgb.DESC).
    Build()
```

## Query Building Patterns

### SELECT Queries

#### Basic Selection
```go
// Select by primary key
user := &User{ID: 123}
query, err := orm.Select().
    Where(qgb.EQ("id")).
    Build()

result, err := query.QueryStruct(ctx, db, user)
```

#### Multiple Conditions
```go
// Using AND conditions
query, err := orm.Select().
    Where(
        qgb.AND(
            qgb.EQ("is_active"),
            qgb.GT("created_at"),
            qgb.IN("category_id"),
        ),
    ).
    Build()

// Using OR conditions
query, err := orm.Select().
    Where(
        qgb.OR(
            qgb.EQ("email"),
            qgb.EQ("phone"),
        ),
    ).
    Build()
```

#### Pagination and Ordering
```go
// Paginated results with ordering
query, err := orm.Select().
    Where(qgb.EQ("is_active")).
    OrderBy("created_at", qgb.DESC).
    Limit(20).
    Offset(40).  // Page 3 with 20 items per page
    Build()

users := []*User{}
err = query.QuerySlice(ctx, db, &users, &User{IsActive: true})
```

#### Field Selection
```go
// Select specific fields only
query, err := orm.Select().
    Fields("id", "email", "name").
    Where(qgb.EQ("is_active")).
    Build()
```

### INSERT Queries

#### Simple Insert
```go
user := &User{
    Email:    "john@example.com",
    Name:     "John Doe",
    IsActive: true,
}

query, err := orm.Insert().
    Returning().  // Get the created record back
    Build()

createdUser, err := query.QueryStruct(ctx, db, user)
// createdUser now contains ID, CreatedAt, UpdatedAt
```

#### Bulk Insert
```go
users := []*User{
    {Email: "user1@example.com", Name: "User 1"},
    {Email: "user2@example.com", Name: "User 2"},
    {Email: "user3@example.com", Name: "User 3"},
}

query, err := orm.Insert().
    Returning().
    Build()

for _, user := range users {
    createdUser, err := query.QueryStruct(ctx, db, user)
    if err != nil {
        log.Printf("Failed to insert user: %v", err)
        continue
    }
    // Process created user
}
```

#### UPSERT (ON CONFLICT)
```go
// Do nothing on conflict
query, err := orm.Insert().
    OnConflict(qgb.DoNothing("email")).
    Returning().
    Build()

// Update on conflict
query, err := orm.Insert().
    OnConflict(qgb.DoUpdate("SET name = EXCLUDED.name, updated_at = NOW()", "email")).
    Returning().
    Build()

// Multiple conflict columns
query, err := orm.Insert().
    OnConflict(qgb.DoUpdate("SET updated_at = NOW()", "email", "phone")).
    Build()
```

### UPDATE Queries

#### Update Specific Fields
```go
user := &User{
    ID:       123,
    Name:     "Updated Name",
    IsActive: false,
}

query, err := orm.Update().
    Fields("name", "is_active", "updated_at").
    Where(qgb.EQ("id")).
    Returning().
    Build()

updatedUser, err := query.QueryStruct(ctx, db, user)
```

#### Conditional Updates
```go
// Update all active users created before a date
query, err := orm.Update().
    Fields("is_active").
    Where(
        qgb.AND(
            qgb.EQ("is_active"),
            qgb.LT("created_at"),
        ),
    ).
    Build()

updateData := &User{
    IsActive:  false,
    CreatedAt: time.Now().AddDate(0, -6, 0), // 6 months ago
}

_, err = query.Exec(ctx, db, updateData)
```

### DELETE Queries

#### Delete by ID
```go
query, err := orm.Delete().
    Where(qgb.EQ("id")).
    Build()

_, err = query.Exec(ctx, db, &User{ID: 123})
```

#### Conditional Delete
```go
// Delete inactive users older than 1 year
query, err := orm.Delete().
    Where(
        qgb.AND(
            qgb.EQ("is_active"),
            qgb.LT("created_at"),
        ),
    ).
    Build()

deleteData := &User{
    IsActive:  false,
    CreatedAt: time.Now().AddDate(-1, 0, 0),
}

result, err := query.Exec(ctx, db, deleteData)
log.Printf("Deleted %d rows", result.RowsAffected())
```

## Advanced Features

### WHERE Clause Operators

#### Comparison Operators
```go
// Equality and inequality
qgb.EQ("field")      // field = @field
qgb.NEQ("field")     // field != @field
qgb.GT("field")      // field > @field
qgb.GTE("field")     // field >= @field
qgb.LT("field")      // field < @field
qgb.LTE("field")     // field <= @field

// Null checks
qgb.ISNULL("field")  // field IS NULL
qgb.NOTNULL("field") // field IS NOT NULL
```

#### Array and Set Operations
```go
// IN clause for multiple values
userIDs := []uint64{1, 2, 3, 4, 5}
query, err := orm.Select().
    Where(qgb.IN("id")).
    Build()

// Execute with slice in struct
users := []*User{}
err = query.QuerySlice(ctx, db, &users, &User{ID: 0}) // Will use userIDs slice

// ANY operator (PostgreSQL specific)
query, err := orm.Select().
    Where(qgb.ANY("tags")).
    Build()
```

#### Raw SQL Clauses
```go
// Use RAW for complex conditions
query, err := orm.Select().
    Where(
        qgb.AND(
            qgb.EQ("is_active"),
            qgb.RAW("created_at > NOW() - INTERVAL '30 days'"),
        ),
    ).
    Build()
```

### Logical Operators

#### Complex Logical Combinations
```go
// Nested logical operations
query, err := orm.Select().
    Where(
        qgb.OR(
            qgb.AND(
                qgb.EQ("is_active"),
                qgb.GT("created_at"),
            ),
            qgb.AND(
                qgb.EQ("is_premium"),
                qgb.GT("last_login"),
            ),
            qgb.IN("id"),
        ),
    ).
    Build()

// NOT operator
query, err := orm.Select().
    Where(
        qgb.NOT(
            qgb.OR(
                qgb.EQ("is_banned"),
                qgb.ISNULL("email"),
            ),
        ),
    ).
    Build()
```

### Query Execution Methods

#### QueryStruct - Single Record
```go
// Get single record
user := &User{ID: 123}
query, err := orm.Select().Where(qgb.EQ("id")).Build()
result, err := query.QueryStruct(ctx, db, user)
if err != nil {
    if err == pgx.ErrNoRows {
        // Handle not found
    }
    return err
}
```

#### QuerySlice - Multiple Records
```go
// Get multiple records
users := []*User{}
query, err := orm.Select().
    Where(qgb.EQ("is_active")).
    OrderBy("created_at", qgb.DESC).
    Build()

err = query.QuerySlice(ctx, db, &users, &User{IsActive: true})
```

#### Exec - No Return Data
```go
// Execute without returning data
query, err := orm.Delete().Where(qgb.EQ("id")).Build()
result, err := query.Exec(ctx, db, &User{ID: 123})
if err != nil {
    return err
}

log.Printf("Deleted %d rows", result.RowsAffected())
```

## Performance Best Practices

### 1. Query Reuse
```go
// Prepare queries once, reuse many times
type UserService struct {
    selectByID    *qgb.Query[User]
    selectActive  *qgb.Query[User]
    insertUser    *qgb.Query[User]
}

func NewUserService(orm *qgb.ORM[User]) (*UserService, error) {
    selectByID, err := orm.Select().Where(qgb.EQ("id")).Build()
    if err != nil {
        return nil, err
    }
    
    selectActive, err := orm.Select().Where(qgb.EQ("is_active")).Build()
    if err != nil {
        return nil, err
    }
    
    insertUser, err := orm.Insert().Returning().Build()
    if err != nil {
        return nil, err
    }
    
    return &UserService{
        selectByID:   selectByID,
        selectActive: selectActive,
        insertUser:   insertUser,
    }, nil
}

func (s *UserService) GetByID(ctx context.Context, db *pgxpool.Pool, id uint64) (*User, error) {
    return s.selectByID.QueryStruct(ctx, db, &User{ID: id})
}
```

### 2. Batch Operations
```go
// Use transactions for batch operations
tx, err := db.Begin(ctx)
if err != nil {
    return err
}
defer tx.Rollback(ctx)

query, err := orm.Insert().Returning().Build()
if err != nil {
    return err
}

createdUsers := make([]*User, 0, len(users))
for _, user := range users {
    created, err := query.QueryStruct(ctx, tx, user)
    if err != nil {
        return err
    }
    createdUsers = append(createdUsers, created)
}

return tx.Commit(ctx)
```

### 3. Field Selection
```go
// Select only needed fields for better performance
query, err := orm.Select().
    Fields("id", "email", "name").  // Only select necessary fields
    Where(qgb.EQ("is_active")).
    Build()
```

## Common Patterns

### 1. Repository Pattern
```go
type UserRepository struct {
    orm *qgb.ORM[User]
    db  *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) (*UserRepository, error) {
    orm, err := qgb.New[User]("users")
    if err != nil {
        return nil, err
    }
    
    return &UserRepository{orm: orm, db: db}, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uint64) (*User, error) {
    query, err := r.orm.Select().Where(qgb.EQ("id")).Build()
    if err != nil {
        return nil, err
    }
    
    return query.QueryStruct(ctx, r.db, &User{ID: id})
}

func (r *UserRepository) Create(ctx context.Context, user *User) (*User, error) {
    query, err := r.orm.Insert().Returning().Build()
    if err != nil {
        return nil, err
    }
    
    return query.QueryStruct(ctx, r.db, user)
}
```

### 2. Service Layer Pattern
```go
type UserService struct {
    repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, email, name string) (*User, error) {
    // Validate input
    if email == "" {
        return nil, errors.New("email is required")
    }
    
    // Check if user exists
    existing, err := s.repo.GetByEmail(ctx, email)
    if err != nil && err != pgx.ErrNoRows {
        return nil, err
    }
    if existing != nil {
        return nil, errors.New("user already exists")
    }
    
    // Create user
    user := &User{
        Email:    email,
        Name:     name,
        IsActive: true,
    }
    
    return s.repo.Create(ctx, user)
}
```

### 3. Pagination Helper
```go
type PaginationParams struct {
    Page     int
    PageSize int
}

func (p PaginationParams) Offset() int {
    return (p.Page - 1) * p.PageSize
}

func (r *UserRepository) GetPaginated(ctx context.Context, params PaginationParams) ([]*User, error) {
    query, err := r.orm.Select().
        Where(qgb.EQ("is_active")).
        OrderBy("created_at", qgb.DESC).
        Limit(params.PageSize).
        Offset(params.Offset()).
        Build()
    
    if err != nil {
        return nil, err
    }
    
    users := []*User{}
    err = query.QuerySlice(ctx, r.db, &users, &User{IsActive: true})
    return users, err
}
```

## Error Handling

### Common Error Patterns
```go
func (r *UserRepository) GetByID(ctx context.Context, id uint64) (*User, error) {
    query, err := r.orm.Select().Where(qgb.EQ("id")).Build()
    if err != nil {
        return nil, fmt.Errorf("failed to build query: %w", err)
    }
    
    user, err := query.QueryStruct(ctx, r.db, &User{ID: id})
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, ErrUserNotFound
        }
        return nil, fmt.Errorf("failed to query user: %w", err)
    }
    
    return user, nil
}

// Define custom errors
var (
    ErrUserNotFound = errors.New("user not found")
    ErrUserExists   = errors.New("user already exists")
)
```

### Constraint Violations
```go
func (r *UserRepository) Create(ctx context.Context, user *User) (*User, error) {
    query, err := r.orm.Insert().Returning().Build()
    if err != nil {
        return nil, err
    }
    
    created, err := query.QueryStruct(ctx, r.db, user)
    if err != nil {
        // Check for unique constraint violations
        if strings.Contains(err.Error(), "unique_violation") {
            return nil, ErrUserExists
        }
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    return created, nil
}
```

## Testing Strategies

### 1. Unit Testing with Mocks
```go
func TestUserService_CreateUser(t *testing.T) {
    repo := &MockUserRepository{}
    service := NewUserService(repo)
    
    // Setup expectations
    repo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, pgx.ErrNoRows)
    repo.On("Create", mock.Anything, mock.AnythingOfType("*User")).Return(&User{
        ID:    1,
        Email: "test@example.com",
        Name:  "Test User",
    }, nil)
    
    // Execute
    user, err := service.CreateUser(context.Background(), "test@example.com", "Test User")
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "test@example.com", user.Email)
    repo.AssertExpectations(t)
}
```

### 2. Integration Testing
```go
func TestUserRepository_Integration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer db.Close()
    
    repo, err := NewUserRepository(db)
    require.NoError(t, err)
    
    // Test creation
    user := &User{
        Email:    "test@example.com",
        Name:     "Test User",
        IsActive: true,
    }
    
    created, err := repo.Create(context.Background(), user)
    require.NoError(t, err)
    assert.NotZero(t, created.ID)
    assert.NotZero(t, created.CreatedAt)
    
    // Test retrieval
    found, err := repo.GetByID(context.Background(), created.ID)
    require.NoError(t, err)
    assert.Equal(t, created.Email, found.Email)
}
```

## Best Practices Summary

1. **Always use struct tags** with proper field mapping
2. **Reuse queries** when possible for better performance
3. **Handle errors appropriately** with custom error types
4. **Use transactions** for batch operations
5. **Select only needed fields** for better performance
6. **Validate inputs** before database operations
7. **Use proper logging** for debugging
8. **Test thoroughly** with both unit and integration tests
9. **Follow repository pattern** for clean architecture
10. **Use connection pooling** for concurrent access

This guide covers the essential patterns for using QGB effectively. The library's design prioritizes performance and type safety while maintaining a clean, intuitive API.