package postgres

var (
	SelectPostById      = "SELECT id, author, title, content, allow_comment, created_at FROM post WHERE id = $1;"
	SelectPostByComment = "SELECT id, author, title, content, allow_comment, created_at FROM post WHERE id = (SELECT path[1] FROM comment WHERE id = $1);"
	SelectFeedPosts     = "SELECT id, author, title, content, allow_comment, created_at FROM post ORDER BY id DESC LIMIT $1 OFFSET $2;"

	SelectFeedReplies = `WITH depth_parent AS (SELECT icount(path) AS depth FROM comment WHERE id = $1) 
SELECT 
	id, author, content,
	path[icount(path)] AS parent, 
	icount(path) - (SELECT depth FROM depth_parent) AS depth,
	created_at 
FROM comment WHERE idx(path[2:], $1) > 0 AND
	($2 < 0 OR icount(path) - (SELECT depth FROM depth_parent) <= $2)
ORDER BY path + intset(id) + (SELECT MAX(id) + 1 FROM comment) DESC LIMIT $3 OFFSET $4;`

	InsertNewPost    = "INSERT INTO post (author, title, content, allow_comment) values($1, $2, $3, $4) RETURNING id;"
	InsertNewComment = "INSERT INTO comment (author, content, path) VALUES ($1, $2, intset($3)) RETURNING id;"
	InsertNewReply   = "INSERT INTO comment (author, content, path) VALUES ($1, $2, (SELECT path FROM comment WHERE id = $3) + $3) RETURNING id;"

	UpdateCommentingPermission = "UPDATE post SET allow_comment = $1 WHERE id = $2 AND allow_comment <> $1;"
)
