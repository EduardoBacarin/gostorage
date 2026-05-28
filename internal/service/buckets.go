package service

import (
	"context"
	"errors"
	"time"

	"github.com/EduardoBacarin/gostorage/internal/helpers"
	"github.com/EduardoBacarin/gostorage/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
	bucketName := helpers.Slugify(name)
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
		ID:          helpers.GenerateSHA256("bucket", bucketName),
		Name:        bucketName,
		AllowPublic: isPublic,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
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
		bucketName := helpers.Slugify(*newName)
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
		currentBucket.AllowPublic = *isPublic // atualiza para o retorno
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

func (s *BucketService) GetBucket(ctx context.Context, bucketID string, userID *string) (*models.Bucket, error) {
	var bucket models.Bucket
	err := s.collection.FindOne(ctx, bson.M{"_id": bucketID}).Decode(&bucket)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("Not Found")
		}
		return nil, err
	}
	return &bucket, nil
}

func (s *BucketService) ListBuckets(ctx context.Context, page, limit int64, sortBy, sortDir, search string, allowedBuckets []string) (*helpers.PaginatedResult[models.Bucket], error) {
	filter := bson.M{}
	if sortBy == "" {
		sortBy = "created_at"
	}
	switch sortBy {
	case "name", "created_at":
	default:
		sortBy = "created_at"
	}

	sortOrder := -1
	if sortDir == "asc" {
		sortOrder = 1
	}

	hasSuperUser := false
	for _, allowed := range allowedBuckets {
		if allowed == "*" {
			hasSuperUser = true
			break
		}
	}

	if !hasSuperUser {
		if len(allowedBuckets) == 0 {
			return helpers.NewPaginatedResult([]models.Bucket{}, 0, page, limit), nil
		}
		filter["_id"] = bson.M{"$in": allowedBuckets}
	}

	if search != "" {
		regexFilter := bson.M{"$regex": search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"name": regexFilter},
			{"_id": regexFilter},
		}
	}

	total, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	pagConfig := helpers.PreparePagination(page, limit)

	findOptions := options.Find()
	findOptions.SetLimit(pagConfig.Limit)
	findOptions.SetSkip(pagConfig.Skip)
	findOptions.SetSort(bson.M{sortBy: sortOrder})

	cursor, err := s.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var buckets []models.Bucket
	if err = cursor.All(ctx, &buckets); err != nil {
		return nil, err
	}

	return helpers.NewPaginatedResult(buckets, total, page, pagConfig.Limit), nil
}

func (s *BucketService) DeleteBucket(ctx context.Context, bucketID string, userID *string) error {
	var bucket models.Bucket
	err := s.collection.FindOne(ctx, bson.M{"_id": bucketID}).Decode(&bucket)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("Not Found")
		}
		return err
	}

	_, err = s.collection.DeleteOne(ctx, bson.M{"_id": bucketID})
	return err
}
