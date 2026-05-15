package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/EduardoBacarin/gostorage/internal/models"
	"github.com/EduardoBacarin/gostorage/internal/storage"
	"github.com/google/uuid"
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

func (s *ObjectService) Upload(ctx context.Context, r io.Reader, bucket, key, owner, contentType string) (*models.ObjectMetadata, error) {
	tempID := uuid.New().String()

	hasher := sha256.New()
	tee := io.TeeReader(r, hasher)

	tempPath, err := s.storage.Save(tempID, tee)
	if err != nil {
		return nil, err
	}

	fileHash := hex.EncodeToString(hasher.Sum(nil))

	var existing models.ObjectMetadata
	err = s.collection.FindOne(ctx, bson.M{"checksum": fileHash}).Decode(&existing)

	if err == nil {
		_ = s.storage.Delete(tempPath)
		return s.createMetadata(ctx, bucket, key, owner, fileHash, contentType, existing.StoragePath)
	}

	finalPath := filepath.Join(filepath.Dir(tempPath), fileHash)

	if err := os.Rename(tempPath, finalPath); err != nil {
		_ = s.storage.Delete(tempPath)
		return nil, err
	}

	return s.createMetadata(ctx, bucket, key, owner, fileHash, contentType, finalPath)
}

func (s *ObjectService) createMetadata(ctx context.Context, bucket, key, owner, hash, contentType, storagePath string) (*models.ObjectMetadata, error) {
	meta := &models.ObjectMetadata{
		ID:          uuid.New().String(),
		OwnerID:     owner,
		Checksum:    hash,
		Bucket:      bucket,
		Key:         key,
		ContentType: contentType,
		StoragePath: storagePath,
		CreatedAt:   time.Now(),
	}

	_, err := s.collection.InsertOne(ctx, meta)
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
