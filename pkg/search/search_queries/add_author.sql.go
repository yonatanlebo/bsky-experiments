// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: add_author.sql

package search_queries

import (
	"context"
)

const addAuthor = `-- name: AddAuthor :exec
INSERT INTO authors (did, handle) VALUES ($1, $2) ON CONFLICT (did) DO UPDATE SET handle = $2
`

type AddAuthorParams struct {
	Did    string `json:"did"`
	Handle string `json:"handle"`
}

func (q *Queries) AddAuthor(ctx context.Context, arg AddAuthorParams) error {
	_, err := q.exec(ctx, q.addAuthorStmt, addAuthor, arg.Did, arg.Handle)
	return err
}
