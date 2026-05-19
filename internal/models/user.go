package models

import "time"

type User struct {
	ID          string    `bson:"_id" json:"id"`
	Email       string    `bson:"email" json:"email"`
	Password    string    `bson:"password" json:"-"`
	Permissions []string  `bson:"permissions"`
	Buckets     []string  `bson:"buckets"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
}
