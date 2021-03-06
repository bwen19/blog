// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: comment.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

const createComment = `-- name: CreateComment :one
INSERT INTO comments (
    parent_id,
    article_id,
    commenter,
    content
) VALUES (
    CASE WHEN $1::bool
        THEN $2::bigint ELSE NULL END,
    $3::bigint,
    $4::varchar,
    $5::varchar
) RETURNING id, parent_id, article_id, commenter, content, comment_at
`

type CreateCommentParams struct {
	SetParentID bool   `json:"set_parent_id"`
	ParentID    int64  `json:"parent_id"`
	ArticleID   int64  `json:"article_id"`
	Commenter   string `json:"commenter"`
	Content     string `json:"content"`
}

func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) (Comment, error) {
	row := q.db.QueryRowContext(ctx, createComment,
		arg.SetParentID,
		arg.ParentID,
		arg.ArticleID,
		arg.Commenter,
		arg.Content,
	)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.ParentID,
		&i.ArticleID,
		&i.Commenter,
		&i.Content,
		&i.CommentAt,
	)
	return i, err
}

const deleteComment = `-- name: DeleteComment :exec
DELETE FROM comments
WHERE id = $1
    AND ($2::bool OR commenter = $3::varchar)
`

type DeleteCommentParams struct {
	ID           int64  `json:"id"`
	AnyCommenter bool   `json:"any_commenter"`
	Commenter    string `json:"commenter"`
}

func (q *Queries) DeleteComment(ctx context.Context, arg DeleteCommentParams) error {
	_, err := q.db.ExecContext(ctx, deleteComment, arg.ID, arg.AnyCommenter, arg.Commenter)
	return err
}

const getComment = `-- name: GetComment :one
SELECT id, parent_id, article_id, commenter, content, comment_at FROM comments
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetComment(ctx context.Context, id int64) (Comment, error) {
	row := q.db.QueryRowContext(ctx, getComment, id)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.ParentID,
		&i.ArticleID,
		&i.Commenter,
		&i.Content,
		&i.CommentAt,
	)
	return i, err
}

const listChildComments = `-- name: ListChildComments :many
SELECT
    id,
    parent_id,
    article_id,
    commenter,
    avatar_src,
    content,
    comment_at
FROM comments AS c
    JOIN users AS u
    ON c.commenter = u.username
WHERE parent_id = ANY($1::bigint[])
`

type ListChildCommentsRow struct {
	ID        int64         `json:"id"`
	ParentID  sql.NullInt64 `json:"parent_id"`
	ArticleID int64         `json:"article_id"`
	Commenter string        `json:"commenter"`
	AvatarSrc string        `json:"avatar_src"`
	Content   string        `json:"content"`
	CommentAt time.Time     `json:"comment_at"`
}

func (q *Queries) ListChildComments(ctx context.Context, commentIds []int64) ([]ListChildCommentsRow, error) {
	rows, err := q.db.QueryContext(ctx, listChildComments, pq.Array(commentIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListChildCommentsRow{}
	for rows.Next() {
		var i ListChildCommentsRow
		if err := rows.Scan(
			&i.ID,
			&i.ParentID,
			&i.ArticleID,
			&i.Commenter,
			&i.AvatarSrc,
			&i.Content,
			&i.CommentAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listCommentsByArticle = `-- name: ListCommentsByArticle :many
SELECT
    id,
    parent_id,
    article_id,
    commenter,
    avatar_src,
    content,
    comment_at
FROM comments AS c
    JOIN users AS u
    ON c.commenter = u.username
WHERE article_id = $3::bigint
    AND parent_id IS NULL
ORDER BY comment_at
LIMIT $1
OFFSET $2
`

type ListCommentsByArticleParams struct {
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
	ArticleID int64 `json:"article_id"`
}

type ListCommentsByArticleRow struct {
	ID        int64         `json:"id"`
	ParentID  sql.NullInt64 `json:"parent_id"`
	ArticleID int64         `json:"article_id"`
	Commenter string        `json:"commenter"`
	AvatarSrc string        `json:"avatar_src"`
	Content   string        `json:"content"`
	CommentAt time.Time     `json:"comment_at"`
}

func (q *Queries) ListCommentsByArticle(ctx context.Context, arg ListCommentsByArticleParams) ([]ListCommentsByArticleRow, error) {
	rows, err := q.db.QueryContext(ctx, listCommentsByArticle, arg.Limit, arg.Offset, arg.ArticleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListCommentsByArticleRow{}
	for rows.Next() {
		var i ListCommentsByArticleRow
		if err := rows.Scan(
			&i.ID,
			&i.ParentID,
			&i.ArticleID,
			&i.Commenter,
			&i.AvatarSrc,
			&i.Content,
			&i.CommentAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
