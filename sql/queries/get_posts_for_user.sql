-- name: GetPostsForUser :many
SELECT 
	posts.id,
	posts.created_at,
	posts.updated_at,
	posts.title,
	posts.url,
	posts.description,
	posts.published_at,
	posts.feed_id
FROM posts
INNER JOIN feed_follows ON feed_follows.feed_id = posts.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC, posts.created_at DESC
LIMIT $2;
