package models

type DownloadQueueRequest struct {
	ID        string `json:"id" bson:"_id"`
	CreatorID int64  `json:"creator_id" bson:"creator_id"`

	SpotifyURL string `json:"spotify_url" bson:"spotify_url"`
	Name       string `json:"name" bson:"name"`
	Active     bool   `json:"active" bson:"active"`
	Errored    bool   `json:"errored" bson:"errored"`

	CreatedAt  int64 `json:"created_at" bson:"created_at"`
	UpdatedAt  int64 `json:"updated_at" bson:"updated_at"`
	SyncCount  int   `json:"sync_count" bson:"sync_count"`
	RetryCount int   `json:"retry_count" bson:"retry_count"`
}

type PlaylistRequest struct {
	ID         string `json:"id" bson:"_id"`
	CreatorID  int64  `json:"creator_id" bson:"creator_id"`
	SpotifyURL string `json:"spotify_url" bson:"spotify_url"`

	Active     bool `json:"active" bson:"active"`
	Errored    bool `json:"errored" bson:"errored"`
	RetryCount int  `json:"retry_count" bson:"retry_count"`

	CreatedAt int64 `json:"created_at" bson:"created_at"`
	UpdatedAt int64 `json:"updated_at" bson:"updated_at"`
}
