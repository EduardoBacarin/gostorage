package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/EduardoBacarin/gostorage/internal/helpers"
	"github.com/EduardoBacarin/gostorage/internal/models"
	"github.com/EduardoBacarin/gostorage/internal/storage"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ObjectService struct {
	storage          storage.Engine
	collection       *mongo.Collection
	bucketCollection *mongo.Collection
}

func NewObjectService(s storage.Engine, db *mongo.Database) *ObjectService {
	return &ObjectService{
		storage:          s,
		collection:       db.Collection("objects"),
		bucketCollection: db.Collection("buckets"),
	}
}

func (s *ObjectService) Upload(ctx context.Context, r io.Reader, bucket, key, owner, contentType string, contentSize int64) (*models.ObjectMetadata, error) {
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
		return s.createMetadata(ctx, bucket, key, owner, fileHash, contentType, contentSize)
	}

	subPath, err := helpers.GenerateDynamicPath(fileHash)
	if err != nil {
		_ = s.storage.Delete(tempPath)
		return nil, err
	}
	baseDir := filepath.Dir(tempPath)
	finalPath := filepath.Join(baseDir, subPath)
	if err := os.MkdirAll(filepath.Dir(finalPath), 0755); err != nil {
		_ = s.storage.Delete(tempPath)
		return nil, err
	}

	if err := os.Rename(tempPath, finalPath); err != nil {
		_ = s.storage.Delete(tempPath)
		return nil, err
	}

	return s.createMetadata(ctx, bucket, key, owner, fileHash, contentType, contentSize)
}

func (s *ObjectService) createMetadata(ctx context.Context, bucket, key, owner, hash, contentType string, size int64) (*models.ObjectMetadata, error) {
	meta := &models.ObjectMetadata{
		ID:          uuid.New().String(),
		OwnerID:     owner,
		Checksum:    hash,
		Bucket:      bucket,
		Key:         key,
		Size:        size,
		ContentType: contentType,
		CreatedAt:   time.Now(),
	}

	_, err := s.collection.InsertOne(ctx, meta)
	if err != nil {
		return nil, err
	}

	return meta, nil
}

func (s *ObjectService) GetObject(ctx context.Context, bucketName string, id string, allowedBuckets []string) (*models.ObjectMetadata, *os.File, error) {

	var bucket models.Bucket
	err := s.bucketCollection.FindOne(ctx, bson.M{"name": bucketName}).Decode(&bucket)

	if err != nil {
		return nil, nil, errors.New("Not found")
	}

	var object models.ObjectMetadata
	err = s.collection.FindOne(ctx, bson.M{"bucket": bucketName, "_id": id}).Decode(&object)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil, errors.New("Not found")
		}
		return nil, nil, err
	}

	if object.IsPublic && bucket.AllowPublic {
		return s.openFileAndReturn(&object)
	}

	hasPermission := false
	for _, allowed := range allowedBuckets {
		if allowed == "*" || allowed == bucketName {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return nil, nil, errors.New("Forbidden")
	}
	return s.openFileAndReturn(&object)
}

func (s *ObjectService) openFileAndReturn(meta *models.ObjectMetadata) (*models.ObjectMetadata, *os.File, error) {
	subPath, err := helpers.GenerateDynamicPath(meta.Checksum)
	if err != nil {
		return nil, nil, err
	}
	finalPath := filepath.Join(s.storage.BaseDir(), subPath)
	file, err := os.Open(finalPath)
	if err != nil {
		return nil, nil, err
	}
	return meta, file, nil
}

func (s *ObjectService) ValidateBucketPermission(ctx context.Context, bucketName string, allowedBuckets []string) error {
	if len(allowedBuckets) == 0 {
		return errors.New("Forbidden")
	}

	hasPermission := false

	for _, allowed := range allowedBuckets {
		if allowed == "*" {
			hasPermission = true
			break
		}
		if allowed == bucketName {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return errors.New("Forbidden")
	}

	return nil
}
