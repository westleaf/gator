-- name: GetNextFeedToFetch :one
SELECT id, url, name, last_fetched_at
FROM feeds
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;
