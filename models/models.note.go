package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       string             `json:"title" validate:"required"`
	Description string             `json:"description" validate:"required"`
	Created_at  int64              `json:"created_at"`
	Updated_at  int64              `json:"updated_At"`
	Note_id     string             `json:"note_id"`
}

type CreateNote struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}
