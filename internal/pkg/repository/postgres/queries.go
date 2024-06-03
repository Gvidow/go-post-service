package postgres

var (
	SelectPostById = "SELECT author, title, content, allow_comment, created_at FROM post WHERE id = $1;"

	SelectFeedReplies = `WITH path_parent_comment AS (SELECT path FROM comment WHERE id = $1), 
	inf AS (SELECT MAX(id) + 1 AS val FROM comment) 
SELECT 
	id, author, content,
	path[icount(path)] AS parent, 
	icount(path) - icount((SELECT path FROM path_parent_comment)) AS depth
	created_at 
FROM comment WHERE path @> (SELECT path FROM path_parent_comment) AND
	($2 < 0 OR icount(path) - icount((SELECT path FROM path_parent_comment)) <= $2)
ORDER BY slice + (SELECT val FROM inf) DESC LIMIT $3 OFFSET $4;`
)
