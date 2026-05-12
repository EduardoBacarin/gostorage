package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"time"

	"github.com/EduardoBacarin/gostorage/internal/models"
	"github.com/EduardoBacarin/gostorage/internal/storage"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ObjectService struct {
	storage    storage.Engine
	collection *mongo.Collection
}

func NewObjectService(s storage.Engine, db *mongo.Database) *ObjectService {
	return &ObjectService{
		storage:    s,
		collection: db.Collection("objects"),
	}
}

func (s *ObjectService) Upload(ctx context.Context, content io.ReadSeeker, bucket, key string) (*models.ObjectMetadata, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, content); err != nil {
		return nil, err
	}
	fileHash := hex.EncodeToString(hash.Sum(nil))

	var existing models.ObjectMetadata
	err := s.collection.FindOne(ctx, bson.M{"_id": fileHash}).Decode(&existing)

	if err == nil {
		return &existing, nil
	}

	content.Seek(0, io.SeekStart)

	path, err := s.storage.Save(fileHash, content)
	if err != nil {
		return nil, err
	}

	meta := &models.ObjectMetadata{
		ID:          fileHash,
		Bucket:      bucket,
		Key:         key,
		StoragePath: path,
		CreatedAt:   time.Now(),
	}

	_, err = s.collection.InsertOne(ctx, meta)
	if err != nil {
		return nil, err
	}

	return meta, nil
}

func (s *ObjectService) GetObject(ctx context.Context, id string) (io.ReadCloser, *models.ObjectMetadata, error) {
	var meta models.ObjectMetadata

	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&meta)
	if err != nil {
		return nil, nil, err
	}

	stream, err := s.storage.Get(meta.StoragePath)
	if err != nil {
		return nil, nil, err
	}

	return stream, &meta, nil
}

func (s *ObjectService) DeleteObject(ctx context.Context, id string) error {
	var meta models.ObjectMetadata

	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&meta)
	if err != nil {
		return err
	}

	_, err = s.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	err = s.storage.Delete(meta.StoragePath)
	if err != nil {
		log.Printf("Database register removed but failed on delete file %s: %v", meta.StoragePath, err)
	}

	return nil
}
