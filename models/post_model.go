package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID     primitive.ObjectID `json:"id,omitempty"`
	Title  string             `json:"title,omitempty" validate:"required"`
	Body   string             `json:"body,omitempty" validate:"required"`
	Author string             `json:"author,omitempty" validate:"required"`
}
