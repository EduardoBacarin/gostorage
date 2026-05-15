package models

import "time"

type Bucket struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Public    bool      `bson:"public" json:"public"`
	OwnerID   string    `bson:"owner_id" json:"owner_id"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
