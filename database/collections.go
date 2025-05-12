package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// downloadQueueRequestCollection returns the download queue request collection
func (d *db) downloadQueueRequestCollection() *mongo.Collection {
	if err := d.conn.Ping(context.Background(), nil); err != nil {
		d.log.Error("failed to ping database. reconnecting.", zap.Error(err))
		d.reconnectToDB()
	}
	return d.conn.Database(d.dbname).Collection(d.downloadRequestCollectionName)
}

func (d *db) playlistsCollection() *mongo.Collection {
	if err := d.conn.Ping(context.Background(), nil); err != nil {
		d.log.Error("failed to ping database. reconnecting.", zap.Error(err))
		d.reconnectToDB()
	}

	return d.conn.Database(d.dbname).Collection(d.playlistRequestCollectionName)
}

func (d *db) indexStatusCollection() *mongo.Collection {
	if err := d.conn.Ping(context.Background(), nil); err != nil {
		d.log.Error("failed to ping database. reconnecting.", zap.Error(err))
		d.reconnectToDB()
	}

	return d.conn.Database(d.dbname).Collection(d.indexStatusCollectionName)
}

// musicFilesCollection returns the music files collection
func (d *db) musicFilesCollection() *mongo.Collection {
	if err := d.conn.Ping(context.Background(), nil); err != nil {
		d.log.Error("failed to ping database. reconnecting.", zap.Error(err))
		d.reconnectToDB()
	}

	return d.conn.Database(d.dbname).Collection(d.musicFilesCollectionName)
}
