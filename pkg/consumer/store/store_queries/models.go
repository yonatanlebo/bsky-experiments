// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package store_queries

import (
	"database/sql"
	"time"

	"github.com/sqlc-dev/pqtype"
)

type Actor struct {
	Did           string         `json:"did"`
	Handle        string         `json:"handle"`
	DisplayName   sql.NullString `json:"display_name"`
	Bio           sql.NullString `json:"bio"`
	HandleValid   bool           `json:"handle_valid"`
	LastValidated sql.NullTime   `json:"last_validated"`
	ProPicCid     sql.NullString `json:"pro_pic_cid"`
	BannerCid     sql.NullString `json:"banner_cid"`
	CreatedAt     sql.NullTime   `json:"created_at"`
	UpdatedAt     sql.NullTime   `json:"updated_at"`
	InsertedAt    time.Time      `json:"inserted_at"`
}

type Block struct {
	ActorDid   string       `json:"actor_did"`
	Rkey       string       `json:"rkey"`
	TargetDid  string       `json:"target_did"`
	CreatedAt  sql.NullTime `json:"created_at"`
	InsertedAt time.Time    `json:"inserted_at"`
}

type Collection struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

type DailySummary struct {
	Date                    time.Time `json:"date"`
	LikesPerDay             int64     `json:"Likes per Day"`
	DailyActiveLikers       int64     `json:"Daily Active Likers"`
	DailyActivePosters      int64     `json:"Daily Active Posters"`
	PostsPerDay             int64     `json:"Posts per Day"`
	PostsWithImagesPerDay   int64     `json:"Posts with Images per Day"`
	ImagesPerDay            int64     `json:"Images per Day"`
	ImagesWithAltTextPerDay int64     `json:"Images with Alt Text per Day"`
	FirstTimePosters        int64     `json:"First Time Posters"`
	FollowsPerDay           int64     `json:"Follows per Day"`
	DailyActiveFollowers    int64     `json:"Daily Active Followers"`
	BlocksPerDay            int64     `json:"Blocks per Day"`
	DailyActiveBlockers     int64     `json:"Daily Active Blockers"`
}

type Event struct {
	ID           int64                 `json:"id"`
	InitiatorDid string                `json:"initiator_did"`
	TargetDid    string                `json:"target_did"`
	EventType    string                `json:"event_type"`
	PostUri      sql.NullString        `json:"post_uri"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
	ExpiredAt    sql.NullTime          `json:"expired_at"`
	ConcludedAt  sql.NullTime          `json:"concluded_at"`
	WindowStart  sql.NullTime          `json:"window_start"`
	WindowEnd    sql.NullTime          `json:"window_end"`
	Results      pqtype.NullRawMessage `json:"results"`
}

type Follow struct {
	ActorDid   string       `json:"actor_did"`
	Rkey       string       `json:"rkey"`
	TargetDid  string       `json:"target_did"`
	CreatedAt  sql.NullTime `json:"created_at"`
	InsertedAt time.Time    `json:"inserted_at"`
}

type FollowerCount struct {
	ActorDid     string    `json:"actor_did"`
	NumFollowers int64     `json:"num_followers"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type FollowerStat struct {
	P25   float64 `json:"p25"`
	P50   float64 `json:"p50"`
	P75   float64 `json:"p75"`
	P90   float64 `json:"p90"`
	P95   float64 `json:"p95"`
	P99   float64 `json:"p99"`
	P995  float64 `json:"p99_5"`
	P997  float64 `json:"p99_7"`
	P999  float64 `json:"p99_9"`
	P9999 float64 `json:"p99_99"`
}

type FollowingCount struct {
	ActorDid     string    `json:"actor_did"`
	NumFollowing int64     `json:"num_following"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Image struct {
	Cid          string         `json:"cid"`
	PostActorDid string         `json:"post_actor_did"`
	PostRkey     string         `json:"post_rkey"`
	AltText      sql.NullString `json:"alt_text"`
	CreatedAt    sql.NullTime   `json:"created_at"`
	InsertedAt   time.Time      `json:"inserted_at"`
}

type Like struct {
	ActorDid   string       `json:"actor_did"`
	Rkey       string       `json:"rkey"`
	Subj       int64        `json:"subj"`
	CreatedAt  sql.NullTime `json:"created_at"`
	InsertedAt time.Time    `json:"inserted_at"`
}

type LikeCount struct {
	SubjectID        int64        `json:"subject_id"`
	NumLikes         int64        `json:"num_likes"`
	UpdatedAt        time.Time    `json:"updated_at"`
	SubjectCreatedAt sql.NullTime `json:"subject_created_at"`
}

type PointAssignment struct {
	ID        int64     `json:"id"`
	EventID   int64     `json:"event_id"`
	ActorDid  string    `json:"actor_did"`
	Points    int32     `json:"points"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Post struct {
	ActorDid           string                `json:"actor_did"`
	Rkey               string                `json:"rkey"`
	Content            sql.NullString        `json:"content"`
	ParentPostActorDid sql.NullString        `json:"parent_post_actor_did"`
	QuotePostActorDid  sql.NullString        `json:"quote_post_actor_did"`
	QuotePostRkey      sql.NullString        `json:"quote_post_rkey"`
	ParentPostRkey     sql.NullString        `json:"parent_post_rkey"`
	RootPostActorDid   sql.NullString        `json:"root_post_actor_did"`
	RootPostRkey       sql.NullString        `json:"root_post_rkey"`
	Facets             pqtype.NullRawMessage `json:"facets"`
	Embed              pqtype.NullRawMessage `json:"embed"`
	Tags               []string              `json:"tags"`
	HasEmbeddedMedia   bool                  `json:"has_embedded_media"`
	CreatedAt          sql.NullTime          `json:"created_at"`
	InsertedAt         time.Time             `json:"inserted_at"`
}

type PostSentiment struct {
	ActorDid    string          `json:"actor_did"`
	Rkey        string          `json:"rkey"`
	InsertedAt  time.Time       `json:"inserted_at"`
	CreatedAt   time.Time       `json:"created_at"`
	ProcessedAt sql.NullTime    `json:"processed_at"`
	Sentiment   sql.NullString  `json:"sentiment"`
	Confidence  sql.NullFloat64 `json:"confidence"`
}

type RecentPostsWithScore struct {
	SubjectID        int64        `json:"subject_id"`
	ActorDid         string       `json:"actor_did"`
	Rkey             string       `json:"rkey"`
	SubjectCreatedAt sql.NullTime `json:"subject_created_at"`
	InsertedAt       time.Time    `json:"inserted_at"`
	Score            float64      `json:"score"`
}

type RepoBackfillStatus struct {
	Repo         string    `json:"repo"`
	LastBackfill time.Time `json:"last_backfill"`
	SeqStarted   int64     `json:"seq_started"`
	State        string    `json:"state"`
}

type Repost struct {
	ActorDid   string       `json:"actor_did"`
	Rkey       string       `json:"rkey"`
	Subj       int64        `json:"subj"`
	CreatedAt  sql.NullTime `json:"created_at"`
	InsertedAt time.Time    `json:"inserted_at"`
}

type RepostCount struct {
	SubjectID  int64     `json:"subject_id"`
	NumReposts int64     `json:"num_reposts"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Subject struct {
	ID       int64  `json:"id"`
	ActorDid string `json:"actor_did"`
	Rkey     string `json:"rkey"`
	Col      int32  `json:"col"`
}
