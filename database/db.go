package database

import (
	"context"

	"github.com/DigitalIndependence/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Database interface {
	GetActiveRequests(ctx context.Context) ([]models.DownloadQueueRequest, error)
	GetActiveRequest(ctx context.Context, url string) (models.DownloadQueueRequest, error)
	CheckIfRequestAlreadySynced(ctx context.Context, url string) (bool, error)
	NewDownloadRequest(ctx context.Context, url, name string, creatorID int64) error
	UpdateActiveRequest(ctx context.Context, request models.DownloadQueueRequest) error

	GetActivePlaylists(ctx context.Context) ([]models.PlaylistRequest, error)
	UpdatePlaylistRequest(ctx context.Context, request models.PlaylistRequest) error
	NewPlaylistRequest(ctx context.Context, url string, creatorID int64) error

	FindMusicFiles(ctx context.Context, artists, titles []string) ([]models.MusicFile, error)
	IndexMusicFile(ctx context.Context, file models.MusicFile) error

	GetIndexStatus(ctx context.Context) (models.IndexStatus, error)
	UpdateIndexStatus(ctx context.Context, status models.IndexStatus) error
}

type db struct {
	conn *mongo.Client
	log  *zap.Logger

	url    string
	dbname string

	// collections
	musicFilesCollectionName      string
	downloadRequestCollectionName string
	playlistRequestCollectionName string
	indexStatusCollectionName     string
}

func NewDatabase(ctx context.Context, log *zap.Logger, url, dbname,
	musicFilesCollection, downloadRequestCollection, indexStatusCollection, playlistRequestCollection string,
) (Database, error) {
	conn, err := mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	if musicFilesCollection == "" || downloadRequestCollection == "" ||
		indexStatusCollection == "" || playlistRequestCollection == "" {
		return nil, ErrEmptyCollectionName
	}

	if dbname == "" {
		return nil, ErrEmptyDBName
	}

	return &db{
		conn: conn,
		log:  log,

		url:    url,
		dbname: dbname,

		// collections
		musicFilesCollectionName:      musicFilesCollection,
		downloadRequestCollectionName: downloadRequestCollection,
		playlistRequestCollectionName: playlistRequestCollection,
		indexStatusCollectionName:     indexStatusCollection,
	}, nil
}

func (d *db) reconnectToDB() error {
	d.conn.Disconnect(context.Background())

	conn, err := mongo.Connect(context.Background(), options.Client().ApplyURI(d.url))
	if err != nil {
		return err
	}

	d.conn = conn
	return nil
}
