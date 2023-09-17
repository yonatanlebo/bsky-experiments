// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: events.sql

package store_queries

import (
	"context"
	"database/sql"

	"github.com/sqlc-dev/pqtype"
)

const addEventPost = `-- name: AddEventPost :exec
UPDATE events
SET post_uri = $2
WHERE id = $1
`

type AddEventPostParams struct {
	ID      int64          `json:"id"`
	PostUri sql.NullString `json:"post_uri"`
}

func (q *Queries) AddEventPost(ctx context.Context, arg AddEventPostParams) error {
	_, err := q.exec(ctx, q.addEventPostStmt, addEventPost, arg.ID, arg.PostUri)
	return err
}

const concludeEvent = `-- name: ConcludeEvent :exec
UPDATE events
SET results = $2,
    concluded_at = $3
WHERE id = $1
`

type ConcludeEventParams struct {
	ID          int64                 `json:"id"`
	Results     pqtype.NullRawMessage `json:"results"`
	ConcludedAt sql.NullTime          `json:"concluded_at"`
}

func (q *Queries) ConcludeEvent(ctx context.Context, arg ConcludeEventParams) error {
	_, err := q.exec(ctx, q.concludeEventStmt, concludeEvent, arg.ID, arg.Results, arg.ConcludedAt)
	return err
}

const confirmEvent = `-- name: ConfirmEvent :exec
UPDATE events
SET window_start = $2,
    window_end = $3,
    expired_at = NULL
WHERE id = $1
`

type ConfirmEventParams struct {
	ID          int64        `json:"id"`
	WindowStart sql.NullTime `json:"window_start"`
	WindowEnd   sql.NullTime `json:"window_end"`
}

func (q *Queries) ConfirmEvent(ctx context.Context, arg ConfirmEventParams) error {
	_, err := q.exec(ctx, q.confirmEventStmt, confirmEvent, arg.ID, arg.WindowStart, arg.WindowEnd)
	return err
}

const createEvent = `-- name: CreateEvent :one
INSERT INTO events (
        initiator_did,
        target_did,
        event_type,
        expired_at
    )
VALUES ($1, $2, $3, $4)
RETURNING id
`

type CreateEventParams struct {
	InitiatorDid string       `json:"initiator_did"`
	TargetDid    string       `json:"target_did"`
	EventType    string       `json:"event_type"`
	ExpiredAt    sql.NullTime `json:"expired_at"`
}

func (q *Queries) CreateEvent(ctx context.Context, arg CreateEventParams) (int64, error) {
	row := q.queryRow(ctx, q.createEventStmt, createEvent,
		arg.InitiatorDid,
		arg.TargetDid,
		arg.EventType,
		arg.ExpiredAt,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const deleteEvent = `-- name: DeleteEvent :exec
DELETE FROM events
WHERE id = $1
`

func (q *Queries) DeleteEvent(ctx context.Context, id int64) error {
	_, err := q.exec(ctx, q.deleteEventStmt, deleteEvent, id)
	return err
}

const getActiveEventsForInitiator = `-- name: GetActiveEventsForInitiator :many
SELECT id, initiator_did, target_did, event_type, post_uri, created_at, updated_at, expired_at, concluded_at, window_start, window_end, results
FROM events
WHERE initiator_did = $1
    AND event_type = $2
    AND (
        window_end > NOW()
        OR expired_at > NOW()
    )
ORDER BY created_at DESC
LIMIT $3 OFFSET $4
`

type GetActiveEventsForInitiatorParams struct {
	InitiatorDid string `json:"initiator_did"`
	EventType    string `json:"event_type"`
	Limit        int32  `json:"limit"`
	Offset       int32  `json:"offset"`
}

func (q *Queries) GetActiveEventsForInitiator(ctx context.Context, arg GetActiveEventsForInitiatorParams) ([]Event, error) {
	rows, err := q.query(ctx, q.getActiveEventsForInitiatorStmt, getActiveEventsForInitiator,
		arg.InitiatorDid,
		arg.EventType,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.InitiatorDid,
			&i.TargetDid,
			&i.EventType,
			&i.PostUri,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ExpiredAt,
			&i.ConcludedAt,
			&i.WindowStart,
			&i.WindowEnd,
			&i.Results,
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

const getActiveEventsForTarget = `-- name: GetActiveEventsForTarget :many
SELECT id, initiator_did, target_did, event_type, post_uri, created_at, updated_at, expired_at, concluded_at, window_start, window_end, results
FROM events
WHERE $1 = target_did
    AND event_type = $2
    AND (
        window_end > NOW()
        OR expired_at < NOW()
    )
ORDER BY created_at DESC
LIMIT $3 OFFSET $4
`

type GetActiveEventsForTargetParams struct {
	TargetDid string `json:"target_did"`
	EventType string `json:"event_type"`
	Limit     int32  `json:"limit"`
	Offset    int32  `json:"offset"`
}

func (q *Queries) GetActiveEventsForTarget(ctx context.Context, arg GetActiveEventsForTargetParams) ([]Event, error) {
	rows, err := q.query(ctx, q.getActiveEventsForTargetStmt, getActiveEventsForTarget,
		arg.TargetDid,
		arg.EventType,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.InitiatorDid,
			&i.TargetDid,
			&i.EventType,
			&i.PostUri,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ExpiredAt,
			&i.ConcludedAt,
			&i.WindowStart,
			&i.WindowEnd,
			&i.Results,
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

const getEvent = `-- name: GetEvent :one
SELECT id, initiator_did, target_did, event_type, post_uri, created_at, updated_at, expired_at, concluded_at, window_start, window_end, results
FROM events
WHERE id = $1
`

func (q *Queries) GetEvent(ctx context.Context, id int64) (Event, error) {
	row := q.queryRow(ctx, q.getEventStmt, getEvent, id)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.InitiatorDid,
		&i.TargetDid,
		&i.EventType,
		&i.PostUri,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiredAt,
		&i.ConcludedAt,
		&i.WindowStart,
		&i.WindowEnd,
		&i.Results,
	)
	return i, err
}

const getEventsForInitiator = `-- name: GetEventsForInitiator :many
SELECT id, initiator_did, target_did, event_type, post_uri, created_at, updated_at, expired_at, concluded_at, window_start, window_end, results
FROM events
WHERE initiator_did = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

type GetEventsForInitiatorParams struct {
	InitiatorDid string `json:"initiator_did"`
	Limit        int32  `json:"limit"`
	Offset       int32  `json:"offset"`
}

func (q *Queries) GetEventsForInitiator(ctx context.Context, arg GetEventsForInitiatorParams) ([]Event, error) {
	rows, err := q.query(ctx, q.getEventsForInitiatorStmt, getEventsForInitiator, arg.InitiatorDid, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.InitiatorDid,
			&i.TargetDid,
			&i.EventType,
			&i.PostUri,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ExpiredAt,
			&i.ConcludedAt,
			&i.WindowStart,
			&i.WindowEnd,
			&i.Results,
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

const getEventsForTarget = `-- name: GetEventsForTarget :many
SELECT id, initiator_did, target_did, event_type, post_uri, created_at, updated_at, expired_at, concluded_at, window_start, window_end, results
FROM events
WHERE $1 = target_did
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

type GetEventsForTargetParams struct {
	TargetDid string `json:"target_did"`
	Limit     int32  `json:"limit"`
	Offset    int32  `json:"offset"`
}

func (q *Queries) GetEventsForTarget(ctx context.Context, arg GetEventsForTargetParams) ([]Event, error) {
	rows, err := q.query(ctx, q.getEventsForTargetStmt, getEventsForTarget, arg.TargetDid, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.InitiatorDid,
			&i.TargetDid,
			&i.EventType,
			&i.PostUri,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ExpiredAt,
			&i.ConcludedAt,
			&i.WindowStart,
			&i.WindowEnd,
			&i.Results,
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

const getEventsToConclude = `-- name: GetEventsToConclude :many
SELECT id, initiator_did, target_did, event_type, post_uri, created_at, updated_at, expired_at, concluded_at, window_start, window_end, results
FROM events
WHERE window_end < NOW()
    AND results IS NULL
    AND event_type = $1
    AND expired_at IS NULL
ORDER BY window_end ASC
LIMIT $2 OFFSET $3
`

type GetEventsToConcludeParams struct {
	EventType string `json:"event_type"`
	Limit     int32  `json:"limit"`
	Offset    int32  `json:"offset"`
}

func (q *Queries) GetEventsToConclude(ctx context.Context, arg GetEventsToConcludeParams) ([]Event, error) {
	rows, err := q.query(ctx, q.getEventsToConcludeStmt, getEventsToConclude, arg.EventType, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.InitiatorDid,
			&i.TargetDid,
			&i.EventType,
			&i.PostUri,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ExpiredAt,
			&i.ConcludedAt,
			&i.WindowStart,
			&i.WindowEnd,
			&i.Results,
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

const getUnconfirmedEvent = `-- name: GetUnconfirmedEvent :one
SELECT id, initiator_did, target_did, event_type, post_uri, created_at, updated_at, expired_at, concluded_at, window_start, window_end, results
FROM events
WHERE post_uri = $1
    AND target_did = $2
    AND window_start IS NULL
    AND window_end IS NULL
    AND results IS NULL
LIMIT 1
`

type GetUnconfirmedEventParams struct {
	PostUri   sql.NullString `json:"post_uri"`
	TargetDid string         `json:"target_did"`
}

func (q *Queries) GetUnconfirmedEvent(ctx context.Context, arg GetUnconfirmedEventParams) (Event, error) {
	row := q.queryRow(ctx, q.getUnconfirmedEventStmt, getUnconfirmedEvent, arg.PostUri, arg.TargetDid)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.InitiatorDid,
		&i.TargetDid,
		&i.EventType,
		&i.PostUri,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiredAt,
		&i.ConcludedAt,
		&i.WindowStart,
		&i.WindowEnd,
		&i.Results,
	)
	return i, err
}
