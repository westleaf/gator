# Gator

A command-line RSS feed aggregator built in Go. Gator allows you to aggregate and browse RSS feeds from multiple sources, storing posts in a local PostgreSQL database.

## Features

- User management (register, login, list users)
- Add and track multiple RSS feeds
- Follow/unfollow specific feeds
- Automatic feed scraping at configurable intervals
- Browse recent posts from your followed feeds
- Persistent storage with PostgreSQL
- Type-safe database queries using sqlc

## Prerequisites

- Go 1.x or higher
- PostgreSQL database running locally
- Database connection string configured

## Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set up your database configuration (ensure your config file contains the database URL)
4. Run database migrations from the `sql/schema/` directory

## Usage

### User Management

```bash
# Register a new user
gator register <username>

# Login as an existing user
gator login <username>

# List all users
gator users

# Reset all users (clear database)
gator reset
```

### Feed Management

```bash
# Add a new feed (requires login)
gator addfeed <feed_name> <feed_url>

# List all feeds
gator feeds

# Follow a feed (requires login)
gator follow <feed_url>

# List feeds you're following (requires login)
gator following

# Unfollow a feed (requires login)
gator unfollow <feed_url>
```

### Reading Posts

```bash
# Start aggregating feeds at intervals (e.g., every 1 minute)
gator agg 1m

# Browse recent posts from your feeds (requires login)
gator browse [limit]
```

## Project Structure

```
gator/
├── internal/
│   ├── config/        # Configuration management
│   └── database/      # Generated sqlc code and queries
├── sql/
│   ├── queries/       # SQL query definitions
│   └── schema/        # Database migrations
├── main.go           # Entry point
├── commands.go       # Command handlers
├── middleware.go     # Authentication middleware
├── rss.go           # RSS feed fetching
└── scrape.go        # Feed scraping logic
```

## Technical Details

- Built with Go's standard library and minimal dependencies
- Uses sqlc for compile-time verified SQL queries
- PostgreSQL for persistent storage
- Context-aware HTTP requests with timeouts
- UUID-based entity identification
- Graceful handling of duplicate posts and malformed data

## Development

The project uses sqlc to generate type-safe Go code from SQL queries. To regenerate database code after modifying queries:

```bash
sqlc generate
```
