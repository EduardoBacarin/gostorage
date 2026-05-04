package models

import "time"

type ObjectMetadata struct {
	ID          string    `bson:"_id" json:"id"`
	Bucket      string    `bson:"bucket" json:"bucket"`
	Key         string    `bson:"key" json:"key"`
	Size        int64     `bson:"size" json:"size"`
	ContentType string    `bson:"content_type" json:"content_type"`
	IsPublic    bool      `bson:"is_public" json:"is_public"`
	OwnerID     string    `bson:"owner_id" json:"owner_id"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	StoragePath string    `bson:"storage_path" json:"storage_path"`
}
