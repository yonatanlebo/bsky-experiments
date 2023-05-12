// Code generated by sqlc. DO NOT EDIT.
// source: get_author_stats.sql

package search_queries

import (
	"context"
)

const getAuthorStats = `-- name: GetAuthorStats :one
WITH postcounts AS (
    SELECT
        author_did,
        COUNT(id) AS num_posts
    FROM
        posts
    GROUP BY
        author_did
),
percentiles AS (
    SELECT
        (percentile_cont(0.25) WITHIN GROUP (ORDER BY num_posts) * 1)::bigint AS p25,
        (percentile_cont(0.50) WITHIN GROUP (ORDER BY num_posts) * 1)::bigint AS p50,
        (percentile_cont(0.75) WITHIN GROUP (ORDER BY num_posts) * 1)::bigint AS p75,
        (percentile_cont(0.90) WITHIN GROUP (ORDER BY num_posts) * 1)::bigint AS p90,
        (percentile_cont(0.95) WITHIN GROUP (ORDER BY num_posts) * 1)::bigint AS p95,
        (percentile_cont(0.99) WITHIN GROUP (ORDER BY num_posts) * 1)::bigint AS p99
    FROM
        postcounts
),
counts AS (
    SELECT
        count(num_posts) AS total,
        SUM(CASE WHEN num_posts > 1 THEN 1 ELSE 0 END) AS gt_1,
        SUM(CASE WHEN num_posts > 5 THEN 1 ELSE 0 END) AS gt_5,
        SUM(CASE WHEN num_posts > 10 THEN 1 ELSE 0 END) AS gt_10,
        SUM(CASE WHEN num_posts > 20 THEN 1 ELSE 0 END) AS gt_20,
        SUM(CASE WHEN num_posts > 100 THEN 1 ELSE 0 END) AS gt_100,
        SUM(CASE WHEN num_posts > 1000 THEN 1 ELSE 0 END)as gt_1000
    FROM
        postcounts
)
SELECT
    total,
    gt_1,
    gt_5,
    gt_10,
    gt_20,
    gt_100,
    gt_1000,
    (SELECT AVG(num_posts) FROM postcounts)::float AS mean,
    p25,
    p50,
    p75,
    p90,
    p95,
    p99
FROM
    counts,
    percentiles
LIMIT 1
`

type GetAuthorStatsRow struct {
	Total  int64   `json:"total"`
	Gt1    int64   `json:"gt_1"`
	Gt5    int64   `json:"gt_5"`
	Gt10   int64   `json:"gt_10"`
	Gt20   int64   `json:"gt_20"`
	Gt100  int64   `json:"gt_100"`
	Gt1000 int64   `json:"gt_1000"`
	Mean   float64 `json:"mean"`
	P25    int64   `json:"p25"`
	P50    int64   `json:"p50"`
	P75    int64   `json:"p75"`
	P90    int64   `json:"p90"`
	P95    int64   `json:"p95"`
	P99    int64   `json:"p99"`
}

func (q *Queries) GetAuthorStats(ctx context.Context) (GetAuthorStatsRow, error) {
	row := q.queryRow(ctx, q.getAuthorStatsStmt, getAuthorStats)
	var i GetAuthorStatsRow
	err := row.Scan(
		&i.Total,
		&i.Gt1,
		&i.Gt5,
		&i.Gt10,
		&i.Gt20,
		&i.Gt100,
		&i.Gt1000,
		&i.Mean,
		&i.P25,
		&i.P50,
		&i.P75,
		&i.P90,
		&i.P95,
		&i.P99,
	)
	return i, err
}
