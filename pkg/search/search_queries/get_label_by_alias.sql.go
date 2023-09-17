// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: get_label_by_alias.sql

package search_queries

import (
	"context"
)

const getLabelByAlias = `-- name: GetLabelByAlias :one
SELECT id, lookup_alias, name
FROM labels
WHERE lookup_alias = $1
`

func (q *Queries) GetLabelByAlias(ctx context.Context, lookupAlias string) (Label, error) {
	row := q.queryRow(ctx, q.getLabelByAliasStmt, getLabelByAlias, lookupAlias)
	var i Label
	err := row.Scan(&i.ID, &i.LookupAlias, &i.Name)
	return i, err
}
