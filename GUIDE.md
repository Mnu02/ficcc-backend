# How `db` Package Works with `main.go`

## Overview

The `db` package manages your Supabase database connection using a **connection pool**. This guide explains how everything connects together.

## Architecture Flow

```
main.go (Application Entry Point)
    ↓
    Loads .env file (godotenv)
    ↓
    Calls db.InitDB() 
    ↓
db/db.go (Database Package)
    ↓
    Creates connection pool (pgxpool)
    ↓
    Stores pool in global variable: DB
    ↓
    Returns to main.go
    ↓
main.go starts HTTP server
    ↓
HTTP Handlers can use db.GetDB() to query database
```

## Step-by-Step Breakdown

### 1. **main.go Startup Sequence**

```go
func main() {
    // Step 1: Load environment variables from .env
    godotenv.Load()  // Reads DATABASE_URL from .env file
    
    // Step 2: Initialize database connection
    db.InitDB()      // Creates connection pool, connects to Supabase
    
    // Step 3: Start HTTP server
    // Now handlers can use db.GetDB() to query database
}
```

### 2. **db.InitDB() - What Happens Inside**

```go
func InitDB() error {
    // 1. Read DATABASE_URL from environment
    dbURL := os.Getenv("DATABASE_URL")
    
    // 2. Parse connection string and configure pool
    config, err := pgxpool.ParseConfig(dbURL)
    
    // 3. Set pool settings:
    //    - Max 25 connections
    //    - Min 5 connections (always ready)
    //    - Connections live for 1 hour max
    //    - Idle connections closed after 30 min
    
    // 4. Create the connection pool
    pool, err := pgxpool.NewWithConfig(ctx, config)
    
    // 5. Test connection with a query
    pool.QueryRow(ctx, "SELECT version()").Scan(&version)
    
    // 6. Store pool in global variable DB
    DB = pool
    
    return nil
}
```

### 3. **Connection Pool Explained**

The `DB` variable is a **connection pool**, not a single connection:

```
┌─────────────────────────────────────┐
│   Connection Pool (DB variable)     │
│                                     │
│  ┌──────┐  ┌──────┐  ┌──────┐      │
│  │Conn 1│  │Conn 2│  │Conn 3│ ...  │  ← Multiple connections
│  └──────┘  └──────┘  └──────┘      │     ready to use
│                                     │
│  Min: 5 connections (always ready)  │
│  Max: 25 connections (scales up)   │
└─────────────────────────────────────┘
         ↓
    Supabase Database
```

**Why use a pool?**
- **Performance**: Reuses connections instead of creating new ones
- **Concurrency**: Multiple requests can use different connections simultaneously
- **Efficiency**: Connections are managed automatically

### 4. **Using the Database in HTTP Handlers**

Here's how you'd use the database connection in your HTTP handlers:

```go
func myHandler(w http.ResponseWriter, r *http.Request) {
    // Get the connection pool
    pool := db.GetDB()
    
    // Create a context (required for database operations)
    ctx := context.Background()
    
    // Query the database
    var result string
    err := pool.QueryRow(ctx, "SELECT name FROM my_table WHERE id = $1", 1).Scan(&result)
    
    // Handle response...
}
```

## Practical Example: Creating a Real Endpoint

Let's create an example endpoint that queries your database:

```go
// In main.go, add this handler:

func usersHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Get database connection pool
    pool := db.GetDB()
    ctx := context.Background()
    
    // Query your users table (replace with your actual table name)
    rows, err := pool.Query(ctx, "SELECT id, name, email FROM users LIMIT 10")
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()
    
    // Build response
    var users []map[string]interface{}
    for rows.Next() {
        var id int
        var name, email string
        if err := rows.Scan(&id, &name, &email); err != nil {
            continue
        }
        users = append(users, map[string]interface{}{
            "id": id,
            "name": name,
            "email": email,
        })
    }
    
    // Send JSON response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

// Register it in main():
http.HandleFunc("/users", usersHandler)
```

## Key Functions in db Package

### `db.InitDB()`
- **Called once** at application startup
- Creates and configures the connection pool
- Tests the connection
- Stores pool in global `DB` variable

### `db.GetDB()`
- **Called whenever** you need to query the database
- Returns the connection pool
- Safe to call from multiple goroutines (handlers)

### `db.CloseDB()`
- **Called once** when application shuts down
- Closes all connections in the pool
- Cleanup function

## Global Variable Pattern

The `DB` variable in `db/db.go` is a **package-level global variable**:

```go
var DB *pgxpool.Pool  // Shared across entire application
```

**Why this pattern?**
- Simple: No need to pass DB around everywhere
- Efficient: One pool shared by all handlers
- Thread-safe: pgxpool handles concurrency

**Alternative pattern** (if you prefer dependency injection):
```go
// Instead of global, pass DB as parameter
func myHandler(db *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
    // Use db here
}
```

## Lifecycle Summary

```
Application Start
    ↓
Load .env → Get DATABASE_URL
    ↓
db.InitDB() → Create pool → Store in DB variable
    ↓
HTTP Server Starts
    ↓
[Request comes in]
    ↓
Handler calls db.GetDB() → Gets pool → Queries database
    ↓
[Response sent]
    ↓
[Repeat for more requests...]
    ↓
Application Shutdown (Ctrl+C)
    ↓
db.CloseDB() → Close all connections
```

## Best Practices

1. **Always use contexts**: Database operations need a `context.Context`
   ```go
   ctx := context.Background()  // Or r.Context() in handlers
   ```

2. **Close rows**: Always `defer rows.Close()` after queries
   ```go
   rows, err := pool.Query(ctx, "...")
   defer rows.Close()  // Important!
   ```

3. **Handle errors**: Always check database errors
   ```go
   if err != nil {
       // Handle error
   }
   ```

4. **Use parameterized queries**: Prevents SQL injection
   ```go
   // Good
   pool.QueryRow(ctx, "SELECT * FROM users WHERE id = $1", userID)
   
   // Bad (SQL injection risk)
   pool.QueryRow(ctx, fmt.Sprintf("SELECT * FROM users WHERE id = %d", userID))
   ```

## Next Steps

1. **Replace example queries** in `db/example.go` with your actual table names
2. **Create handlers** in `main.go` that use `db.GetDB()` to query your tables
3. **Add error handling** and proper JSON responses
4. **Test endpoints** with curl or Postman

