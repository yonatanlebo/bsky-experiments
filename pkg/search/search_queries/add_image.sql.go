// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: add_image.sql

package search_queries

import (
	"context"
	"database/sql"
	"time"

	"github.com/sqlc-dev/pqtype"
)

const addImage = `-- name: AddImage :exec
INSERT INTO images (cid, post_id, author_did, alt_text, mime_type, fullsize_url, thumbnail_url, created_at, cv_completed, cv_run_at, cv_classes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`

type AddImageParams struct {
	Cid          string                `json:"cid"`
	PostID       string                `json:"post_id"`
	AuthorDid    string                `json:"author_did"`
	AltText      sql.NullString        `json:"alt_text"`
	MimeType     string                `json:"mime_type"`
	FullsizeUrl  string                `json:"fullsize_url"`
	ThumbnailUrl string                `json:"thumbnail_url"`
	CreatedAt    time.Time             `json:"created_at"`
	CvCompleted  bool                  `json:"cv_completed"`
	CvRunAt      sql.NullTime          `json:"cv_run_at"`
	CvClasses    pqtype.NullRawMessage `json:"cv_classes"`
}

func (q *Queries) AddImage(ctx context.Context, arg AddImageParams) error {
	_, err := q.exec(ctx, q.addImageStmt, addImage,
		arg.Cid,
		arg.PostID,
		arg.AuthorDid,
		arg.AltText,
		arg.MimeType,
		arg.FullsizeUrl,
		arg.ThumbnailUrl,
		arg.CreatedAt,
		arg.CvCompleted,
		arg.CvRunAt,
		arg.CvClasses,
	)
	return err
}
