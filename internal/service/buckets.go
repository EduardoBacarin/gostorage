package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/EduardoBacarin/gostorage/internal/helpers"
	"github.com/EduardoBacarin/gostorage/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BucketService struct {
	collection *mongo.Collection
}

func NewBucketService(db *mongo.Database) *BucketService {
	return &BucketService{
		collection: db.Collection("buckets"),
	}
}

func (s *BucketService) CreateBucket(ctx context.Context, name string, isPublic bool, ownerID string) (*models.Bucket, error) {
	bucketName := strings.ToLower(strings.TrimSpace(name))
	if bucketName == "" {
		return nil, errors.New("Bucket name cannot be empty")
	}

	var existing models.Bucket
	err := s.collection.FindOne(ctx, bson.M{"name": bucketName}).Decode(&existing)

	if err == nil {
		return nil, errors.New("Bucket already exists")
	}

	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	newBucket := &models.Bucket{
		ID:        helpers.GenerateSHA256("bucket", bucketName),
		Name:      bucketName,
		Public:    isPublic,
		OwnerID:   ownerID,
		CreatedAt: time.Now(),
	}

	_, err = s.collection.InsertOne(ctx, newBucket)
	if err != nil {
		return nil, err
	}

	return newBucket, nil
}

func (s *BucketService) UpdateBucket(ctx context.Context, bucketID string, newName *string, isPublic *bool, ownerID string) (*models.Bucket, error) {
	var currentBucket models.Bucket
	err := s.collection.FindOne(ctx, bson.M{"_id": bucketID}).Decode(&currentBucket)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("Not Found")
		}
		return nil, err
	}

	updateFields := bson.M{}

	if newName != nil {
		bucketName := strings.ToLower(strings.TrimSpace(*newName))
		if bucketName == "" {
			return nil, errors.New("Bucket name cannot be empty")
		}

		if currentBucket.Name != bucketName {
			var existing models.Bucket
			err = s.collection.FindOne(ctx, bson.M{"name": bucketName}).Decode(&existing)
			if err == nil {
				return nil, errors.New("Bucket already exists")
			}
			if !errors.Is(err, mongo.ErrNoDocuments) {
				return nil, err
			}
		}
		updateFields["name"] = bucketName
		currentBucket.Name = bucketName // atualiza para o retorno
	}

	if isPublic != nil {
		updateFields["public"] = *isPublic
		currentBucket.Public = *isPublic // atualiza para o retorno
	}

	if len(updateFields) == 0 {
		return &currentBucket, nil
	}

	filter := bson.M{"_id": bucketID}
	_, err = s.collection.UpdateOne(ctx, filter, bson.M{"$set": updateFields})
	if err != nil {
		return nil, err
	}

	return &currentBucket, nil
}
