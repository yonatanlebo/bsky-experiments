// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: get_top_posters.sql

package search_queries

import (
	"context"
)

const getTopPosters = `-- name: GetTopPosters :many
SELECT post_count, handle, author_did
FROM top_posters
LIMIT $1
`

func (q *Queries) GetTopPosters(ctx context.Context, limit int32) ([]TopPoster, error) {
	rows, err := q.query(ctx, q.getTopPostersStmt, getTopPosters, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TopPoster
	for rows.Next() {
		var i TopPoster
		if err := rows.Scan(&i.PostCount, &i.Handle, &i.AuthorDid); err != nil {
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
