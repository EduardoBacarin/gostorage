package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
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

func NewObjectService(s storage.Engine, col *mongo.Collection) *ObjectService {
	return &ObjectService{
		storage:    s,
		collection: col,
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
